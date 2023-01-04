package voice

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

// OpusPacketHeaderSize is the size of the opus packet header.
const OpusPacketHeaderSize = 12

// ErrDecryptionFailed is returned when the packet decryption fails.
var ErrDecryptionFailed = errors.New("decryption failed")

var (
	_ io.Reader      = (UDPConn)(nil)
	_ io.ReadCloser  = (UDPConn)(nil)
	_ io.Writer      = (UDPConn)(nil)
	_ io.WriteCloser = (UDPConn)(nil)
	_ net.Conn       = (UDPConn)(nil)
)

type (
	// UDPConnCreateFunc is a function that creates a UDPConn.
	UDPConnCreateFunc func(opts ...UDPConnConfigOpt) UDPConn

	// UDPConn represents a UDP connection to discord voice servers. It is used to send/receive voice packets to/from discord.
	// It implements the io.Reader, io.Writer and io.Closer interface.
	UDPConn interface {
		// LocalAddr returns the local network address, if known.
		LocalAddr() net.Addr

		// RemoteAddr returns the remote network address, if known.
		RemoteAddr() net.Addr

		// SetSecretKey sets the secret key used to encrypt packets.
		SetSecretKey(secretKey [32]byte)

		SetDeadline(t time.Time) error

		// SetReadDeadline sets the read deadline for the UDPConn connection.
		SetReadDeadline(t time.Time) error

		// SetWriteDeadline sets the write deadline for the UDPConn connection.
		SetWriteDeadline(t time.Time) error

		// Open opens the UDPConn connection.
		Open(ctx context.Context, ip string, port int, ssrc uint32) (string, int, error)

		// Close closes the UDPConn connection.
		Close() error

		// Read reads a packet from the UDPConn connection. This implements the io.Reader interface.
		Read(p []byte) (int, error)

		// ReadPacket reads a packet from the UDPConn connection.
		ReadPacket() (*Packet, error)

		// Write writes a packet to the UDPConn connection. This implements the io.Writer interface.
		Write(p []byte) (int, error)
	}

	// Packet is a voice packet received from discord.
	Packet struct {
		// Sequence is the sequence number of the packet.
		Sequence uint16
		// Timestamp is the timestamp of the packet.
		Timestamp uint32
		// SSRC is the users SSRC of the packet.
		SSRC uint32
		// Opus is the actual opus data of the packet.
		Opus []byte
	}
)

// NewUDPConn creates a new voice UDPConn with the given configuration.
func NewUDPConn(opts ...UDPConnConfigOpt) UDPConn {
	config := DefaultUDPConnConfig()
	config.Apply(opts)

	return &udpConnImpl{
		config:        config,
		receiveBuffer: make([]byte, 1400),
	}
}

type udpConnImpl struct {
	config UDPConnConfig

	ip   string
	port int
	ssrc uint32

	conn   net.Conn
	connMu sync.Mutex

	packet    [12]byte
	secretKey [32]byte

	sequence  uint16
	timestamp uint32
	nonce     [24]byte

	receiveNonce  [24]byte
	receiveBuffer []byte
}

func (u *udpConnImpl) LocalAddr() net.Addr {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.LocalAddr()
}

func (u *udpConnImpl) RemoteAddr() net.Addr {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.RemoteAddr()
}

func (u *udpConnImpl) SetSecretKey(secretKey [32]byte) {
	u.secretKey = secretKey
}

func (u *udpConnImpl) SetDeadline(t time.Time) error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.SetDeadline(t)
}

func (u *udpConnImpl) SetReadDeadline(t time.Time) error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.SetReadDeadline(t)
}

func (u *udpConnImpl) SetWriteDeadline(t time.Time) error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.SetWriteDeadline(t)
}

func (u *udpConnImpl) Open(ctx context.Context, ip string, port int, ssrc uint32) (string, int, error) {
	u.ip = ip
	u.port = port
	u.ssrc = ssrc

	u.connMu.Lock()
	defer u.connMu.Unlock()
	host := net.JoinHostPort(u.ip, strconv.Itoa(u.port))
	u.config.Logger.Debugf("Opening UDPConn connection to: %s\n", host)
	var err error
	u.conn, err = u.config.Dialer.DialContext(ctx, "udp", host)
	if err != nil {
		return "", 0, fmt.Errorf("failed to open UDPConn connection: %w", err)
	}

	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, u.ssrc)
	if _, err = u.conn.Write(sb); err != nil {
		return "", 0, fmt.Errorf("failed to write ssrc to UDPConn connection: %w", err)
	}

	rb := make([]byte, 70)
	if _, err = u.conn.Read(rb); err != nil {
		return "", 0, fmt.Errorf("failed to read ip discovery from UDPConn connection: %w", err)
	}

	ourAddress := rb[4:68]
	ourPort := binary.BigEndian.Uint16(rb[68:70])

	u.packet = [12]byte{
		0: 0x80, // Version + Flags
		1: 0x78, // Payload Type
		// [2:4] // Sequence
		// [4:8] // Timestamp
	}

	binary.BigEndian.PutUint32(u.packet[8:12], u.ssrc) // SSRC

	return strings.Replace(string(ourAddress), "\x00", "", -1), int(ourPort), nil
}

func (u *udpConnImpl) Write(p []byte) (int, error) {
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

func (u *udpConnImpl) Read(p []byte) (n int, err error) {
	packet, err := u.ReadPacket()
	if err != nil {
		return 0, err
	}
	return copy(p, packet.Opus), nil
}

func (u *udpConnImpl) ReadPacket() (*Packet, error) {
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

func (u *udpConnImpl) Close() error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.Close()
}
