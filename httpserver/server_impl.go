package httpserver

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
)

var _ Server = (*serverImpl)(nil)

// New creates a new Server with the given publicKey eventHandlerFunc and ConfigOpt(s)
func New(publicKey string, eventHandlerFunc EventHandlerFunc, opts ...ConfigOpt) Server {
	cfg := defaultConfig()
	cfg.apply(opts)

	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		cfg.Logger.Debug("error while decoding hex string", slog.Any("err", err))
	}

	return &serverImpl{
		config:           cfg,
		publicKey:        hexDecodedKey,
		eventHandlerFunc: eventHandlerFunc,
	}
}

type serverImpl struct {
	config           config
	publicKey        PublicKey
	eventHandlerFunc EventHandlerFunc
}

func (s *serverImpl) Start() {
	s.config.ServeMux.Handle(s.config.URL, HandleInteraction(s.publicKey, s.config.Logger, s.eventHandlerFunc))
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
