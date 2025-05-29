package httpserver

import "crypto/ed25519"

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
