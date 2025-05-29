package httpserver

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// Server is used for receiving Discord's interactions via Outgoing Webhooks
type Server interface {
	// Start starts the Server
	Start()

	// Close closes the Server
	Close(ctx context.Context)
}

var _ Server = (*serverImpl)(nil)

// New creates a new Server with the given publicKey eventHandlerFunc and ConfigOpt(s)
func New(publicKey string, eventHandlerFunc EventHandlerFunc, opts ...ConfigOpt) (Server, error) {
	cfg := defaultConfig()
	cfg.apply(opts)

	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding public key: %w", err)
	}

	return &serverImpl{
		config:           cfg,
		publicKey:        hexDecodedKey,
		eventHandlerFunc: eventHandlerFunc,
		verifier:         cfg.Verifier,
	}, nil
}

type serverImpl struct {
	config           config
	publicKey        PublicKey
	eventHandlerFunc EventHandlerFunc
	verifier         Verifier
}

func (s *serverImpl) Start() {
	s.config.ServeMux.Handle(s.config.URL, HandleInteraction(s.verifier, s.publicKey, s.config.Logger, s.eventHandlerFunc))
	s.config.HTTPServer.Addr = s.config.Address
	s.config.HTTPServer.Handler = s.config.ServeMux

	go func() {
		var err error
		if s.config.CertFile != "" && s.config.KeyFile != "" {
			err = s.config.HTTPServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
		} else {
			err = s.config.HTTPServer.ListenAndServe()
		}
		if !errors.Is(err, http.ErrServerClosed) {
			s.config.Logger.Error("error while running http server", slog.Any("err", err))
		}
	}()
}

func (s *serverImpl) Close(ctx context.Context) {
	_ = s.config.HTTPServer.Shutdown(ctx)
}
