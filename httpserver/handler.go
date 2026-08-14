package httpserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/json/v2"

	"github.com/disgoorg/disgo/discord"
)

type (
	// EventHandlerFunc is used to handle events from Discord's Outgoing Webhooks
	EventHandlerFunc func(responseFunc RespondFunc, event EventInteractionCreate)

	// RespondFunc is used to respond to Discord's Outgoing Webhooks
	RespondFunc func(response discord.InteractionResponse) error
)

// EventInteractionCreate is the event payload when an interaction is created via Discord's Outgoing Webhooks
type EventInteractionCreate struct {
	discord.Interaction
}

func (e *EventInteractionCreate) UnmarshalJSON(data []byte) error {
	interaction, err := discord.UnmarshalInteraction(data)
	if err != nil {
		return err
	}
	e.Interaction = interaction
	return nil
}

func (e EventInteractionCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Interaction)
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

type replyStatus int

const (
	replyStatusWaiting replyStatus = iota
	replyStatusReplied
	replyStatusTimedOut
)

// HandleInteraction handles an interaction from Discord's Outgoing Webhooks. It verifies and parses the interaction and then calls the passed EventHandlerFunc.
func HandleInteraction(verifier Verifier, publicKey PublicKey, logger *slog.Logger, handleFunc EventHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// we only read & keep bodies around for logging if debug logging is enabled
		debug := logger.Enabled(r.Context(), slog.LevelDebug)

		if ok := VerifyRequest(verifier, r, publicKey); !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			if debug {
				data, _ := io.ReadAll(r.Body)
				logger.DebugContext(r.Context(), "received http interaction with invalid signature", slog.String("body", string(data)))
			}
			return
		}

		defer func() {
			_ = r.Body.Close()
		}()

		rqBody := io.Reader(r.Body)
		if debug {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				logger.ErrorContext(r.Context(), "error while reading interaction body", slog.Any("err", err))
			}
			logger.DebugContext(r.Context(), "received http interaction", slog.String("body", string(data)))
			rqBody = bytes.NewReader(data)
		}

		var v EventInteractionCreate
		if err := json.NewDecoder(rqBody).Decode(&v); err != nil {
			logger.Error("error while decoding interaction", slog.Any("err", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// these channels are used to communicate between the http handler and where the interaction is responded to
		responseChannel := make(chan discord.InteractionResponse, 1)
		defer close(responseChannel)
		errorChannel := make(chan error, 1)
		defer close(errorChannel)

		// status of this interaction with a mutex to ensure usage between multiple goroutines
		var (
			status replyStatus
			mu     sync.Mutex
		)

		// send interaction to our handler
		go handleFunc(func(response discord.InteractionResponse) error {
			mu.Lock()
			defer mu.Unlock()

			if status == replyStatusTimedOut {
				return discord.ErrInteractionExpired
			}

			if status == replyStatusReplied {
				return discord.ErrInteractionAlreadyReplied
			}

			status = replyStatusReplied
			responseChannel <- response
			// wait if we get any error while processing the response
			return <-errorChannel
		}, v)

		var (
			body any
			err  error
		)

		// wait for the interaction to be responded to or to time out after 3s
		ctx, cancel := context.WithTimeout(context.Background(), 3100*time.Millisecond)
		defer cancel()
		select {
		case response := <-responseChannel:
			// if we only acknowledge the interaction, we don't need to send a response body
			// we just need to send a 202 Accepted status
			if response.Type == discord.InteractionResponseTypeAcknowledge {
				w.WriteHeader(http.StatusAccepted)
				return
			}

			if body, err = response.ToBody(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				errorChannel <- err
				return
			}

		case <-ctx.Done():
			mu.Lock()
			defer mu.Unlock()
			status = replyStatusTimedOut

			logger.Debug("interaction timed out")
			http.Error(w, "Interaction Timed Out", http.StatusRequestTimeout)
			errorChannel <- ctx.Err()
			return
		}

		var rsBody bytes.Buffer
		rsWriter := io.Writer(w)
		if debug {
			rsWriter = io.MultiWriter(w, &rsBody)
		}

		if multiPart, ok := body.(*discord.MultipartBuffer); ok {
			w.Header().Set("Content-Type", multiPart.ContentType)
			_, err = io.Copy(rsWriter, multiPart.Buffer)
		} else {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(rsWriter).Encode(body)
		}
		if err != nil {
			errorChannel <- err
			return
		}

		if debug {
			logger.DebugContext(r.Context(), "response to http interaction", slog.String("body", rsBody.String()))
		}
	}
}
