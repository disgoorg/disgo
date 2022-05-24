package httpserver

import "crypto/ed25519"

var (
	// Verify reports whether sig is a valid signature of message by publicKey. It
	Verify = func(publicKey PublicKey, message []byte, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}

	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = ed25519.SignatureSize
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey = []byte
