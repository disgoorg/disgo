package voice

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

// NewEncrypter creates a new Encrypter based on the given encryption mode and secret key.
// If the encryption mode is not supported, an error is returned.
func NewEncrypter(encryptionMode EncryptionMode, secretKey []byte) (Encrypter, error) {
	switch encryptionMode {
	case EncryptionModeNone:
		return NewNoopEncrypter(), nil
	case EncryptionModeAEADAES256GCMRTPSize, EncryptionModeAEADXChaCha20Poly1305RTPSize:
		return NewAEADEncrypter(encryptionMode, secretKey)

	default:
		return nil, fmt.Errorf("unknown encryption mode: %s", encryptionMode)
	}
}

// Encrypter is used to encrypt RTP packets before sending them to Discord.
//
// The header is the 12 byte RTP header.
// The data is the opus encoded audio data.
//
// The returned byte slice is the encrypted packet ready to be sent to Discord.
//
// [NoopEncrypter] does not encrypt the data and is used for testing purposes only.
// [AEADEncrypter] is the required aead_xchacha20_poly1305_rtpsize encryption mode by Discord.
// [AEADAES256GCMRTPSize] is the preferred aead_aes256_gcm_rtpsize encryption mode by Discord.
// See https://discord.com/developers/docs/topics/voice-connections#transport-encryption-and-sending-voice for more information.
type Encrypter interface {
	// Encrypt encrypts the given RTP header and opus data and returns the encrypted packet.
	Encrypt(header [RTPHeaderSize]byte, data []byte) []byte

	// Decrypt decrypts the given packet and returns the RTP header and opus data.
	Decrypt(rtpHeaderSize int, packet []byte) ([]byte, error)
}

// NewNoopEncrypter creates a new NoopEncrypter.
func NewNoopEncrypter() *NoopEncrypter {
	return &NoopEncrypter{
		buf:    make([]byte, RTPHeaderSize, MaxOpusFrameSize+RTPHeaderSize),
		recBuf: make([]byte, 0, MaxOpusFrameSize+RTPHeaderSize),
	}
}

// NoopEncrypter is used to not encrypt RTP packets. This is only for testing purposes.
// Do not use this in production as Discord requires encryption.
type NoopEncrypter struct {
	buf    []byte
	recBuf []byte
}

func (n *NoopEncrypter) Encrypt(header [RTPHeaderSize]byte, data []byte) []byte {
	n.buf = n.buf[:RTPHeaderSize]

	copy(n.buf, header[:])
	n.buf = append(n.buf, data...)

	return n.buf
}

func (n *NoopEncrypter) Decrypt(rtpHeaderSize int, packet []byte) ([]byte, error) {
	n.recBuf = n.recBuf[:0]
	copy(n.recBuf, packet)

	return n.recBuf[rtpHeaderSize:], nil
}

// NewAEADEncrypter creates a new AEADEncrypter with the given encryption mode and secret key.
func NewAEADEncrypter(encryptionMode EncryptionMode, secretKey []byte) (*AEADEncrypter, error) {
	var aead cipher.AEAD

	switch encryptionMode {
	case EncryptionModeAEADAES256GCMRTPSize:
		block, err := aes.NewCipher(secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create AES cipher: %w", err)
		}

		c, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
		}
		aead = c
	case EncryptionModeAEADXChaCha20Poly1305RTPSize:
		c, err := chacha20poly1305.NewX(secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create XChaCha20-Poly1305 cipher: %w", err)
		}
		aead = c
	default:
		return nil, fmt.Errorf("unknown encryption mode: %s", encryptionMode)
	}

	maxFrameSize := MaxOpusFrameSize + RTPHeaderSize + aead.NonceSize() + aead.Overhead()

	return &AEADEncrypter{
		cipher:   aead,
		buf:      make([]byte, RTPHeaderSize, maxFrameSize),
		nonce:    make([]byte, aead.NonceSize()),
		seq:      0,
		recBuf:   make([]byte, 0, maxFrameSize),
		recNonce: make([]byte, aead.NonceSize()),
	}, nil
}

// AEADEncrypter is used to encrypt RTP packets using AEAD ciphers.
type AEADEncrypter struct {
	cipher cipher.AEAD
	buf    []byte
	nonce  []byte
	seq    uint32

	recBuf   []byte
	recNonce []byte
}

func (a *AEADEncrypter) Encrypt(header [RTPHeaderSize]byte, data []byte) []byte {
	a.buf = a.buf[:RTPHeaderSize]

	binary.LittleEndian.PutUint32(a.nonce, a.seq)
	a.seq++

	copy(a.buf, header[:])
	a.buf = a.cipher.Seal(a.buf, a.nonce, data, header[:])
	a.buf = append(a.buf, a.nonce[:4]...)

	return a.buf
}

func (a *AEADEncrypter) Decrypt(rtpHeaderSize int, packet []byte) ([]byte, error) {
	a.recBuf = a.recBuf[:0]

	copy(a.recNonce, packet[len(packet)-4:])

	var err error
	a.recBuf, err = a.cipher.Open(a.recBuf, a.recNonce, packet[rtpHeaderSize:len(packet)-4], packet[:rtpHeaderSize])
	if err != nil {
		return nil, err
	}

	return a.recBuf, nil
}
