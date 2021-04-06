package api

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// WebhookServer is used for receiving an Interaction over http
type WebhookServer interface {
	Disgo() Disgo
	PublicKey() ed25519.PublicKey
	ListenURL() string
	Router() *mux.Router
	Start() error
	Close()
}

// Verify implements the verification side of the discord interactions api signing algorithm, as documented here:
// https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
// Credit: https://github.com/bsdlp/discord-interactions-go/blob/main/interactions/verify.go
func Verify(r *http.Request, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
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
			log.Errorf("error while closing request body: %s", err)
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

	return ed25519.Verify(key, msg.Bytes(), sig)
}
