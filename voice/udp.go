package voice

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrInvalidPacket    = errors.New("invalid packet")
)

type (
	UDPConnCreateFunc func(ip string, port int, ssrc uint32, opts ...UDPConnConfigOpt) *UDPConn
)

var _ io.Reader = (*UDPConn)(nil)

func NewUDPConn(ip string, port int, ssrc uint32, opts ...UDPConnConfigOpt) *UDPConn {
	config := DefaultUDPConnConfig()
	config.Apply(opts)

	return &UDPConn{
		config: *config,
		ip:     ip,
		port:   port,
		ssrc:   ssrc,
	}
}

type UDPConn struct {
	config UDPConnConfig

	ip   string
	port int

	conn   net.Conn
	connMu sync.Mutex
	ssrc   uint32

	packet    [12]byte
	secretKey [32]byte

	sequence  uint16
	timestamp uint32
	nonce     [24]byte

	receiveNonce  [24]byte
	receiveBuffer []byte
	receiveOpus   []byte
}

func (c *UDPConn) HandleGatewayMessageSessionDescription(data GatewayMessageDataSessionDescription) {
	c.secretKey = data.SecretKey
}

func (c *UDPConn) Open(ctx context.Context) (string, int, error) {
	conn, err := c.config.Dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		return "", 0, err
	}
	c.conn = conn

	// ip discovery
	ssrcBuffer := [70]byte{
		0x1, 0x2,
	}
	binary.BigEndian.PutUint16(ssrcBuffer[2:4], 70)
	binary.BigEndian.PutUint32(ssrcBuffer[4:8], c.ssrc)

	_, err = conn.Write(ssrcBuffer[:])
	if err != nil {
		return "", 0, err
	}

	var addressBuffer [70]byte
	_, err = io.ReadFull(conn, addressBuffer[:])
	if err != nil {
		return "", 0, err
	}

	addressBody := addressBuffer[4:68]

	nullPos := bytes.Index(addressBody, []byte{'\x00'})
	if nullPos < 0 {
		return "", 0, errors.New("UDP IP discovery did not contain a null terminator")
	}

	address := addressBody[:nullPos]
	port := binary.LittleEndian.Uint16(addressBody[68:70])

	c.packet = [12]byte{
		0: 0x80, // Version + Flags
		1: 0x78, // Payload Type
		// [2:4] // Sequence
		// [4:8] // Timestamp
	}

	binary.BigEndian.PutUint32(c.packet[8:12], c.ssrc)

	return string(address), int(port), nil
}

func (c *UDPConn) Write(frameLength time.Duration) io.Writer {
	return WriterFunc(func(b []byte) (int, error) {
		// Write a new sequence.
		binary.BigEndian.PutUint16(c.packet[2:4], c.sequence)
		c.sequence++

		binary.BigEndian.PutUint32(c.packet[4:8], c.timestamp)
		c.timestamp += uint32(frameLength.Milliseconds())

		// Copy the first 12 bytes from the packet into the nonce.
		copy(c.nonce[:12], c.packet[:])

		return c.conn.Write(secretbox.Seal(c.packet[:12], b, &c.nonce, &c.secretKey))
	})
}

func (c *UDPConn) Read(b []byte) (n int, err error) {
	i, err := c.conn.Read(c.receiveBuffer)
	if err != nil {
		return 0, err
	}

	if i < 12 || (c.receiveBuffer[0] != 0x80 && c.receiveBuffer[0] != 0x90) {
		return 0, ErrInvalidPacket
	}

	copy(c.receiveNonce[:], c.receiveBuffer[0:12])

	var ok bool
	c.receiveOpus, ok = secretbox.Open(c.receiveBuffer[12:], c.receiveBuffer[12:], &c.receiveNonce, &c.secretKey)
	if !ok {
		return 0, ErrDecryptionFailed
	}

	isExtension := c.receiveBuffer[0]&0x10 == 0x10
	isMarker := c.receiveBuffer[1]&0x80 != 0x0

	if isExtension && !isMarker {
		extLen := binary.BigEndian.Uint16(c.receiveOpus[2:4])
		shift := 4 + 4*int(extLen)

		if len(c.receiveOpus) > shift {
			c.receiveOpus = c.receiveOpus[shift:]
		}
	}

	return copy(b, c.receiveOpus), nil
}

func (c *UDPConn) Close() {
	_ = c.conn.Close()
}

type WriterFunc func([]byte) (int, error)

func (f WriterFunc) Write(b []byte) (int, error) {
	return f(b)
}
