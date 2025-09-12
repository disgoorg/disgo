package httpinteraction

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/json/v2"

	"github.com/disgoorg/disgo/discord"
)

type (
	// InteractionHandlerFunc is used to handle events from Discord's Outgoing Webhooks
	InteractionHandlerFunc func(respond RespondFunc, event EventInteractionCreate)

	// RespondFunc is used to respond to Discord's Outgoing Webhooks
	RespondFunc func(response discord.InteractionResponse) error
)

type HTTPHandler interface {
	Handle(pattern string, handler http.Handler)
}

// New creates a new Server with the given publicKey eventHandlerFunc and ConfigOpt(s)
func New(handler HTTPHandler, publicKey string, interactionHandler InteractionHandlerFunc, opts ...ConfigOpt) error {
	cfg := defaultConfig()
	cfg.apply(opts)

	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return fmt.Errorf("error decoding public key: %w", err)
	}

	handler.Handle(cfg.Endpoint, HandleInteraction(cfg.Verifier, hexDecodedKey, cfg.Logger, interactionHandler))

	return nil
}

type replyStatus int

const (
	replyStatusWaiting replyStatus = iota
	replyStatusReplied
	replyStatusTimedOut
)

// HandleInteraction handles an interaction from Discord's Outgoing Webhooks. It verifies and parses the interaction and then calls the passed EventHandlerFunc.
func HandleInteraction(verifier KeyVerifier, publicKey PublicKey, logger *slog.Logger, handleFunc InteractionHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		if ok := VerifyRequest(verifier, r, publicKey); !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			data, _ := io.ReadAll(r.Body)
			logger.Debug("received http interaction with invalid signature", slog.String("body", string(data)))
			return
		}

		buff := new(bytes.Buffer)
		rqData, _ := io.ReadAll(io.TeeReader(r.Body, buff))
		logger.Debug("received http interaction", slog.String("body", string(rqData)))

		var v EventInteractionCreate
		if err := json.NewDecoder(buff).Decode(&v); err != nil {
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

		rsBody := &bytes.Buffer{}
		multiWriter := io.MultiWriter(w, rsBody)

		if multiPart, ok := body.(*discord.MultipartBuffer); ok {
			w.Header().Set("Content-Type", multiPart.ContentType)
			_, err = io.Copy(multiWriter, multiPart.Buffer)
		} else {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(multiWriter).Encode(body)
		}
		if err != nil {
			errorChannel <- err
			return
		}

		rsData, _ := io.ReadAll(rsBody)
		logger.Debug("response to http interaction", slog.String("body", string(rsData)))
	}
}
