package httpserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

type EventHandlerFunc func(responseChannel chan<- discord.InteractionResponse, payload io.Reader)

// Server is used for receiving an Interaction over httpserver
type Server interface {
	Logger() log.Logger
	PublicKey() PublicKey
	Handle(c chan discord.InteractionResponse, payload io.Reader)

	Start() error
	Close(ctx context.Context)
}

// VerifyRequest implements the verification side of the discord interactions api signing algorithm, as documented here: https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
// Credit: https://github.com/bsdlp/discord-interactions-go/blob/main/interactions/verify.go
func VerifyRequest(logger log.Logger, r *http.Request, key PublicKey) bool {
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
		err = r.Body.Close()
		if err != nil {
			logger.Error("error while closing request body: ", err)
		}
	}()
	var body bytes.Buffer

	defer func() {
		r.Body = ioutil.NopCloser(&body)
	}()

	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return Verify(key, msg.Bytes(), sig)
}
