package httpgateway

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/disgoorg/disgo/discord"
)

type (
	EventHandlerFunc func(ack func(), message Message)

	// InteractionHandlerFunc is used to handle events from Discord's Outgoing Webhooks
	InteractionHandlerFunc func(responseFunc RespondFunc, event EventInteractionCreate)

	// RespondFunc is used to respond to Discord's Outgoing Webhooks
	RespondFunc func(response discord.InteractionResponse) error
)

type Gateway interface {
	Start()

	Close(ctx context.Context)
}

var _ Gateway = (*gatewayImpl)(nil)

// New creates a new Server with the given publicKey eventHandlerFunc and ConfigOpt(s)
func New(publicKey string, interactionHandler InteractionHandlerFunc, eventHandler EventHandlerFunc, opts ...ConfigOpt) (Gateway, error) {
	cfg := defaultConfig()
	cfg.apply(opts)

	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding public key: %w", err)
	}

	return &gatewayImpl{
		config:             cfg,
		publicKey:          hexDecodedKey,
		interactionHandler: interactionHandler,
		eventHandler:       eventHandler,
		verifier:           cfg.Verifier,
	}, nil
}

type gatewayImpl struct {
	config             config
	publicKey          PublicKey
	interactionHandler InteractionHandlerFunc
	eventHandler       EventHandlerFunc
	verifier           Verifier
}

func (s *gatewayImpl) Start() {
	s.config.ServeMux.Handle(s.config.InteractionURL, HandleInteraction(s.verifier, s.publicKey, s.config.Logger, s.interactionHandler))
	s.config.ServeMux.Handle(s.config.EventURL, HandleEvent(s.verifier, s.publicKey, s.config.Logger, s.eventHandler))
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

func (s *gatewayImpl) Close(ctx context.Context) {
	_ = s.config.HTTPServer.Shutdown(ctx)
}
