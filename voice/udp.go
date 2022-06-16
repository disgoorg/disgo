package voice

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

const OpusPacketHeaderSize = 12

var ErrDecryptionFailed = errors.New("decryption failed")

type UDPCreateFunc func(ip string, port int, ssrc uint32, opts ...UDPConfigOpt) UDP

type UDP interface {
	SetSecretKey(secretKey [32]byte)
	Open(ctx context.Context) (string, int, error)
	Write(p []byte) (int, error)
	Read(p []byte) (int, error)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	ReadPacket() (*Packet, error)
	Close()
}

func NewUDP(ip string, port int, ssrc uint32, opts ...UDPConfigOpt) UDP {
	config := DefaultUDPConfig()
	config.Apply(opts)

	return &udpImpl{
		config:        *config,
		ip:            ip,
		port:          port,
		ssrc:          ssrc,
		receiveBuffer: make([]byte, 1400),
	}
}

type udpImpl struct {
	config UDPConfig

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
}

func (u *udpImpl) SetSecretKey(secretKey [32]byte) {
	u.secretKey = secretKey
}

func (u *udpImpl) SetReadDeadline(t time.Time) error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.SetReadDeadline(t)
}

func (u *udpImpl) SetWriteDeadline(t time.Time) error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.SetWriteDeadline(t)
}

func (u *udpImpl) Open(ctx context.Context) (string, int, error) {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	u.config.Logger.Debugf("Opening UDP connection to: %s:%d\n", u.ip, u.port)
	var err error
	u.conn, err = u.config.Dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", u.ip, u.port))
	if err != nil {
		return "", 0, fmt.Errorf("failed to open UDP connection: %w", err)
	}

	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, u.ssrc)
	if _, err = u.conn.Write(sb); err != nil {
		return "", 0, fmt.Errorf("failed to write ssrc to UDP connection: %w", err)
	}

	rb := make([]byte, 70)
	if _, err = u.conn.Read(rb); err != nil {
		return "", 0, fmt.Errorf("failed to read ip discovery from UDP connection: %w", err)
	}

	address := rb[4:68]
	port := binary.BigEndian.Uint16(rb[68:70])

	u.packet = [12]byte{
		0: 0x80, // Version + Flags
		1: 0x78, // Payload Type
		// [2:4] // Sequence
		// [4:8] // Timestamp
	}

	binary.BigEndian.PutUint32(u.packet[8:12], u.ssrc) // SSRC

	return strings.Replace(string(address), "\x00", "", -1), int(port), nil
}

func (u *udpImpl) Write(p []byte) (int, error) {
	binary.BigEndian.PutUint16(u.packet[2:4], u.sequence)
	u.sequence++

	binary.BigEndian.PutUint32(u.packet[4:8], u.timestamp)
	u.timestamp += 960

	// Copy the first 12 bytes from the packet into the nonce.
	copy(u.nonce[:12], u.packet[:])

	u.connMu.Lock()
	conn := u.conn
	u.connMu.Unlock()
	if _, err := conn.Write(secretbox.Seal(u.packet[:], p, &u.nonce, &u.secretKey)); err != nil {
		return 0, fmt.Errorf("failed to write packet: %w", err)
	}
	return len(p), nil
}

func (u *udpImpl) Read(p []byte) (n int, err error) {
	packet, err := u.ReadPacket()
	if err != nil {
		return 0, err
	}
	return copy(p, packet.Opus), nil
}

func (u *udpImpl) ReadPacket() (*Packet, error) {
	u.connMu.Lock()
	conn := u.conn
	u.connMu.Unlock()

	for {
		i, err := conn.Read(u.receiveBuffer)
		if err != nil {
			return nil, fmt.Errorf("failed to read packet: %w", err)
		}
		if i < OpusPacketHeaderSize || (u.receiveBuffer[0] != 0x80 && u.receiveBuffer[0] != 0x90) || (u.receiveBuffer[1] != 0x78 && u.receiveBuffer[1] != 0x80) {
			continue
		}

		copy(u.receiveNonce[:], u.receiveBuffer[0:OpusPacketHeaderSize])

		opus, ok := secretbox.Open(nil, u.receiveBuffer[OpusPacketHeaderSize:i], &u.receiveNonce, &u.secretKey)
		if !ok {
			return nil, ErrDecryptionFailed
		}

		isExtension := u.receiveBuffer[0]&0x10 == 0x10
		isMarker := u.receiveBuffer[1]&0x80 != 0x0

		if isExtension && !isMarker {
			extLen := binary.BigEndian.Uint16(opus[2:4])
			shift := 4 + 4*int(extLen)

			if len(opus) > shift {
				opus = opus[shift:]
			}
		}
		return &Packet{
			Sequence:  binary.BigEndian.Uint16(u.receiveBuffer[2:4]),
			Timestamp: binary.BigEndian.Uint32(u.receiveBuffer[4:8]),
			SSRC:      binary.BigEndian.Uint32(u.receiveBuffer[8:12]),
			Opus:      opus,
		}, nil
	}
}

func (u *udpImpl) Close() {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	_ = u.conn.Close()
}

type Packet struct {
	Sequence  uint16
	Timestamp uint32
	SSRC      uint32
	Opus      []byte
}
