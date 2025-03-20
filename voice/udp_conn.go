package voice

import (
	"context"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// OpusPacketHeaderSize is the size of the opus packet header.
	OpusPacketHeaderSize = 12

	// UDPTimeout is the timeout for UDP connections.
	UDPTimeout = 30 * time.Second
)

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
		SetSecretKey(encryptionMode EncryptionMode, secretKey []byte)

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
	config.Logger = config.Logger.With(slog.String("name", "voice_conn_udp_conn"))

	return &udpConnImpl{
		config:        config,
		receiveBuffer: make([]byte, 1400),
	}
}

type udpConnImpl struct {
	config UDPConnConfig

	conn   net.Conn
	connMu sync.Mutex

	cipher cipher.AEAD

	packet [12]byte

	sequence       uint16
	i              uint32
	timestamp      uint32
	nonce          [24]byte
	associatedData [12]byte

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

func (u *udpConnImpl) SetSecretKey(encryptionMode EncryptionMode, secretKey []byte) {
	var (
		c   cipher.AEAD
		err error
	)
	switch encryptionMode {
	case EncryptionModeAEADXChaCha20Poly1305RTPSize:
		c, err = chacha20poly1305.NewX(secretKey)
	default:
		u.config.Logger.Error("unknown encryption mode", slog.String("mode", string(encryptionMode)))
		return
	}
	if err != nil {
		u.config.Logger.Error("failed to create cipher", slog.Any("err", err))
		return
	}
	u.cipher = c
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
	u.connMu.Lock()
	defer u.connMu.Unlock()
	host := net.JoinHostPort(ip, strconv.Itoa(port))
	u.config.Logger.Debug("Opening UDPConn connection", slog.String("host", host))
	var err error
	u.conn, err = u.config.Dialer.DialContext(ctx, "udp", host)
	if err != nil {
		return "", 0, fmt.Errorf("failed to open UDPConn connection: %w", err)
	}

	// see payload here https://discord.com/developers/docs/topics/voice-connections#ip-discovery
	sb := make([]byte, 74)
	binary.BigEndian.PutUint16(sb[:2], 1)      // 1 = send
	binary.BigEndian.PutUint16(sb[2:4], 70)    // 70 = length
	binary.BigEndian.PutUint32(sb[4:74], ssrc) // ssrc

	if err = u.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return "", 0, fmt.Errorf("failed to set write deadline on UDPConn connection: %w", err)
	}
	defer func() {
		_ = u.conn.SetWriteDeadline(time.Time{})
	}()
	if _, err = u.conn.Write(sb); err != nil {
		return "", 0, fmt.Errorf("failed to write ssrc to UDPConn connection: %w", err)
	}

	rb := make([]byte, 74)
	if err = u.conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return "", 0, fmt.Errorf("failed to set read deadline on UDPConn connection: %w", err)
	}
	defer func() {
		_ = u.conn.SetReadDeadline(time.Time{})
	}()
	if _, err = u.conn.Read(rb); err != nil {
		return "", 0, fmt.Errorf("failed to read ip discovery from UDPConn connection: %w", err)
	}

	if binary.BigEndian.Uint16(rb[0:2]) != 2 {
		return "", 0, fmt.Errorf("invalid ip discovery response")
	}

	size := binary.BigEndian.Uint16(rb[2:4])
	if size != 70 {
		return "", 0, fmt.Errorf("invalid ip discovery response size")
	}

	returnedSSRC := binary.BigEndian.Uint32(rb[4:8])   // ssrc
	ourAddress := strings.TrimSpace(string(rb[8:72]))  // our ip
	ourPort := int(binary.BigEndian.Uint16(rb[72:74])) // our port

	if returnedSSRC != ssrc {
		return "", 0, fmt.Errorf("invalid ssrc in ip discovery response")
	}

	u.packet = [12]byte{
		0: 0x80, // Version + Flags
		1: 0x78, // Payload Type
		// [2:4] // Sequence
		// [4:8] // Timestamp
	}

	binary.BigEndian.PutUint32(u.packet[8:12], ssrc) // SSRC

	return ourAddress, ourPort, nil
}

func (u *udpConnImpl) Write(p []byte) (int, error) {
	binary.BigEndian.PutUint16(u.packet[2:4], u.sequence)
	u.sequence++

	u.i++
	binary.BigEndian.PutUint32(u.associatedData[:], u.i)

	binary.BigEndian.PutUint32(u.packet[4:8], u.timestamp)
	u.timestamp += 960

	// Copy the first 12 bytes from the packet into the nonce.
	copy(u.nonce[:12], u.packet[:])

	u.connMu.Lock()
	conn := u.conn
	u.connMu.Unlock()

	//secretbox.Seal(u.packet[:], p, &u.nonce, &u.secretKey)
	if _, err := conn.Write(u.cipher.Seal(u.packet[12:], u.nonce[:], p, u.associatedData[:])); err != nil {
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

		opus, ok := secretbox.Open(nil, u.receiveBuffer[OpusPacketHeaderSize:i], &u.receiveNonce, nil)
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
