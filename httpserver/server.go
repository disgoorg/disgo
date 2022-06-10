package httpserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

type (
	// EventHandlerFunc is used to handle events from Discord's Outgoing Webhooks
	EventHandlerFunc func(responseFunc RespondFunc, payload io.Reader)

	// RespondFunc is used to respond to Discord's Outgoing Webhooks
	RespondFunc func(response discord.InteractionResponse) error
)

// Server is used for receiving Discord's interactions via Outgoing Webhooks
type Server interface {
	// Logger returns the logger used by the Server
	Logger() log.Logger

	// PublicKey returns the public key used by the Server
	PublicKey() PublicKey

	// Start starts the Server
	Start()

	// Close closes the Server
	Close(ctx context.Context)

	// Handle passes a payload to the Server for processing
	Handle(respondFunc RespondFunc, payload io.Reader)
}

// VerifyRequest implements the verification side of the discord interactions api signing algorithm, as documented here: https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
// Credit: https://github.com/bsdlp/discord-interactions-go/blob/main/interactions/verify.go
func VerifyRequest(r *http.Request, key PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != SignatureSize || sig[63]&224 != 0 {
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

	return Verify(key, msg.Bytes(), sig)
}
