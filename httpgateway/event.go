package httpgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

func HandleEvent(verifier Verifier, publicKey PublicKey, logger *slog.Logger, handleFunc EventHandlerFunc) http.HandlerFunc {
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
		logger.Debug("received http event", slog.String("body", string(rqData)))

		var v Message
		if err := json.NewDecoder(buff).Decode(&v); err != nil {
			logger.Error("error while decoding event", slog.Any("err", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ackChan := make(chan struct{}, 1)
		defer close(ackChan)

		go handleFunc(sync.OnceFunc(func() {
			ackChan <- struct{}{}
		}), v)

		ctx, cancel := context.WithTimeout(context.Background(), 3100*time.Millisecond)
		defer cancel()

		select {
		case <-ackChan:
			w.WriteHeader(http.StatusNoContent)
		case <-ctx.Done():
			logger.Debug("event timed out")
			http.Error(w, "Event Timed Out", http.StatusRequestTimeout)
		}
	}
}
