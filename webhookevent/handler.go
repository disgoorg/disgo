package webhookevent

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

	"github.com/disgoorg/disgo/httpinteraction"
)

type EventHandlerFunc func(ack func(), message Message)

type HTTPHandler interface {
	Handle(pattern string, handler http.Handler)
}

// New creates a new Server with the given publicKey eventHandlerFunc and ConfigOpt(s)
func New(handler HTTPHandler, publicKey string, eventHandler EventHandlerFunc, opts ...ConfigOpt) error {
	cfg := defaultConfig()
	cfg.apply(opts)

	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return fmt.Errorf("error decoding public key: %w", err)
	}

	handler.Handle(cfg.Endpoint, HandleEvent(cfg.Verifier, hexDecodedKey, cfg.Logger, cfg.EnableRawEvents, eventHandler))

	return nil
}

func HandleEvent(verifier httpinteraction.KeyVerifier, publicKey httpinteraction.PublicKey, logger *slog.Logger, raw bool, handle EventHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		if ok := httpinteraction.VerifyRequest(verifier, r, publicKey); !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			data, _ := io.ReadAll(r.Body)
			logger.Debug("received http interaction with invalid signature", slog.String("body", string(data)))
			return
		}

		buff := new(bytes.Buffer)
		rqData, _ := io.ReadAll(io.TeeReader(r.Body, buff))
		logger.Debug("received http event", slog.String("body", string(rqData)))

		var message Message
		if err := json.NewDecoder(buff).Decode(&message); err != nil {
			logger.Error("error while decoding event", slog.Any("err", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ackChan := make(chan struct{}, 1)
		defer close(ackChan)

		ack := sync.OnceFunc(func() {
			ackChan <- struct{}{}
		})

		if raw && message.Type == MessageTypeEvent {
			event := message.Event.(Event)

			go handle(ack, Message{
				Version:       message.Version,
				ApplicationID: message.ApplicationID,
				Type:          message.Type,
				Event: Event{
					Type:      event.Type,
					Timestamp: event.Timestamp,
					Data: EventDataRaw{
						EventType: event.Type,
						Payload:   bytes.NewReader(event.RawData),
					},
					RawData: event.RawData,
				},
				RawEvent: message.RawEvent,
			})
		}

		go handle(ack, message)

		ctx, cancel := context.WithTimeout(context.Background(), 3100*time.Millisecond)
		defer cancel()

		select {
		case <-ackChan:
			w.WriteHeader(http.StatusNoContent)
		case <-ctx.Done():
			logger.Debug("event timed out")
			http.Error(w, "event Timed Out", http.StatusRequestTimeout)
		}
	}
}
