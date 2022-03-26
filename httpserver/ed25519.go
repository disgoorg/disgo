package httpserver

import "crypto/ed25519"

var (
	Verify = func(publicKey PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}
	SignatureSize = ed25519.SignatureSize
)

type (
	PublicKey = []byte
)
