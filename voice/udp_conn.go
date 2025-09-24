package voice

import (
	"context"
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
)

const (
	// RTPHeaderSize is the size of the opus packet header.
	RTPHeaderSize = 12

	// RTPVersionPadExtend is the first byte of the RTP header.
	// Bit index 0 and 1 represent the RTP Protocol version used. Discord uses the latest RTP protocol version, 2.
	// Bit index 2 represents whether we pad. Opus uses an internal padding system, so RTP padding is not used.
	// Bit index 3 represents if we use extensions.
	// Bit index 4 to 7 represent the CC or CSRC count. CSRC is Combined SSRC.
	RTPVersionPadExtend = 0x80

	// RTPPayloadType is Discord's RTP Profile Payload type.
	// I've yet to find actual documentation on what the bits inside this value represent.
	RTPPayloadType = 0x78

	// MaxOpusFrameSize is the max size of an opus frame in bytes.
	MaxOpusFrameSize = 1400

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
		SetSecretKey(encryptionMode EncryptionMode, secretKey []byte) error

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
		Type byte
		// Sequence is the sequence number of the packet.
		Sequence uint16
		// Timestamp is the timestamp of the packet.
		Timestamp uint32
		// SSRC is the users SSRC of the packet.
		SSRC         uint32
		HasExtension bool
		ExtensionID  int
		Extension    []byte
		CSRC         []uint32
		HeaderSize   int
		// Opus is the actual opus data of the packet.
		Opus []byte
	}
)

// NewUDPConn creates a new voice UDPConn with the given configuration.
func NewUDPConn(opts ...UDPConnConfigOpt) UDPConn {
	cfg := defaultUDPConnConfig()
	cfg.apply(opts)

	return &udpConnImpl{
		config:        cfg,
		receiveBuffer: make([]byte, 1400),
	}
}

type udpConnImpl struct {
	config udpConnConfig

	conn   net.Conn
	connMu sync.Mutex

	encrypter Encrypter

	header    [12]byte
	sequence  uint16
	timestamp uint32

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

func (u *udpConnImpl) SetSecretKey(encryptionMode EncryptionMode, secretKey []byte) error {
	e, err := NewEncrypter(encryptionMode, secretKey)
	if err != nil {
		return fmt.Errorf("failed to create encrypter: %w", err)
	}

	u.encrypter = e
	return nil
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

	u.header[0] = RTPVersionPadExtend // Version + Flags
	u.header[1] = RTPPayloadType      // Payload Type
	// [2:4]  // Sequence
	// [4:8]  // Timestamp
	// [8:12] // SSRC

	binary.BigEndian.PutUint32(u.header[8:], ssrc) // SSRC

	return ourAddress, ourPort, nil
}

func (u *udpConnImpl) Write(p []byte) (int, error) {
	u.connMu.Lock()
	conn := u.conn
	u.connMu.Unlock()

	binary.BigEndian.PutUint16(u.header[2:4], u.sequence)
	u.sequence++

	binary.BigEndian.PutUint32(u.header[4:8], u.timestamp)
	u.timestamp += OpusFrameSize

	if _, err := conn.Write(u.encrypter.Encrypt(u.header, p)); err != nil {
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
		n, err := conn.Read(u.receiveBuffer)
		if err != nil {
			return nil, fmt.Errorf("failed to read packet: %w", err)
		}

		packetType := u.receiveBuffer[1]
		if packetType != RTPPayloadType {
			// ignore non-voice packets
			continue
		}

		hasPadding := (u.receiveBuffer[0] & 0x20) != 0
		if hasPadding {
			paddingLen := int(u.receiveBuffer[n-1])
			u.receiveBuffer = u.receiveBuffer[:n-paddingLen]
			n -= paddingLen
		}

		p := Packet{
			Type:         packetType,
			Sequence:     binary.BigEndian.Uint16(u.receiveBuffer[2:4]),
			Timestamp:    binary.BigEndian.Uint32(u.receiveBuffer[4:8]),
			SSRC:         binary.BigEndian.Uint32(u.receiveBuffer[8:RTPHeaderSize]),
			HasExtension: (u.receiveBuffer[0] & 0x10) != 0,
			ExtensionID:  0,
			Extension:    nil,
			CSRC:         nil,
			HeaderSize:   0,
			Opus:         nil,
		}

		cc := int(u.receiveBuffer[0] & 0x0F)
		p.CSRC = make([]uint32, cc)
		offset := RTPHeaderSize
		for i := range cc {
			p.CSRC[i] = binary.BigEndian.Uint32(u.receiveBuffer[offset : offset+4])
			offset += 4
		}

		var extensionLen int
		if p.HasExtension {
			p.ExtensionID = int(binary.BigEndian.Uint16(u.receiveBuffer[offset : offset+2]))
			offset += 2
			extensionLen = int(binary.BigEndian.Uint16(u.receiveBuffer[offset : offset+2]))
			offset += 2
		}

		p.HeaderSize = offset

		decrypted, err := u.encrypter.Decrypt(p.HeaderSize, u.receiveBuffer[:n])
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt packet: %w", err)
		}

		var decryptedOffset int
		if p.HasExtension {
			extensionLen *= 4 // extension length is in 32-bit words
			p.Extension = decrypted[decryptedOffset : decryptedOffset+extensionLen]
			decryptedOffset += extensionLen
		}

		return &Packet{
			Type:         RTPPayloadType,
			Sequence:     p.Sequence,
			Timestamp:    p.Timestamp,
			SSRC:         p.SSRC,
			Extension:    p.Extension,
			HasExtension: p.HasExtension,
			CSRC:         nil,
			HeaderSize:   RTPHeaderSize,
			Opus:         decrypted[decryptedOffset:],
		}, nil
	}
}

func (u *udpConnImpl) Close() error {
	u.connMu.Lock()
	defer u.connMu.Unlock()
	return u.conn.Close()
}
