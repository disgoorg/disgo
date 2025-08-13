package httpgateway

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"net/http"
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey = []byte

// Verifier is used to verify Ed25519 signatures.
type Verifier interface {
	// Verify verifies the signature of the message using the public key.
	// It returns true if the signature is valid, false otherwise.
	Verify(publicKey PublicKey, message []byte, sig []byte) bool

	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize() int
}

// DefaultVerifier is the default implementation of the Verifier interface.
type DefaultVerifier struct{}

func (DefaultVerifier) Verify(publicKey PublicKey, message []byte, sig []byte) bool {
	return ed25519.Verify(publicKey, message, sig)
}

func (DefaultVerifier) SignatureSize() int {
	return ed25519.SignatureSize
}

// VerifyRequest implements the verification side of the discord interactions api signing algorithm, as documented here: https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
// Credit: https://github.com/bsdlp/discord-interactions-go/blob/main/interactions/verify.go
func VerifyRequest(verifier Verifier, r *http.Request, key PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != verifier.SignatureSize() || sig[63]&224 != 0 {
		return false
	}

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer func() {
		_ = r.Body.Close()
	}()
	var body bytes.Buffer

	defer func() {
		r.Body = io.NopCloser(&body)
	}()

	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return verifier.Verify(key, msg.Bytes(), sig)
}
