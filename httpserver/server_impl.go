package httpserver

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

var _ Server = (*ServerImpl)(nil)

func New(logger log.Logger, config Config, eventHandler EventHandlerFunc) Server {
	hexDecodedKey, err := hex.DecodeString(config.PublicKey)
	if err != nil {
		logger.Errorf("error while decoding hex string: %s", err)
	}
	mux := http.NewServeMux()
	mux.Handle(config.URL, &WebhookInteractionHandler{})
	return &ServerImpl{
		logger:           logger,
		publicKey:        hexDecodedKey,
		config:           config,
		eventHandlerFunc: eventHandler,
		server: &http.Server{
			Addr:    config.Port,
			Handler: mux,
		},
	}
}

// ServerImpl is used in Disgo's webhook server for interactions
type ServerImpl struct {
	logger           log.Logger
	config           Config
	publicKey        ed25519.PublicKey
	server           *http.Server
	eventHandlerFunc EventHandlerFunc
}

func (w *ServerImpl) Logger() log.Logger {
	return w.logger
}

// Config returns the Config
func (w *ServerImpl) Config() Config {
	return w.config
}

// Start makes the ServerImpl listen on the specified port and handle requests
func (w *ServerImpl) Start() {
	go func() {
		if err := w.server.ListenAndServe(); err != nil {
			w.logger.Errorf("error starting http server: %s", err)
		}
	}()
}

// Close shuts down the ServerImpl
func (w *ServerImpl) Close() {
	if err := w.server.Close(); err != nil {
		w.Logger().Errorf("error while shutting down http server: %s", err)
	}
}

type WebhookInteractionHandler struct {
	server *ServerImpl
}

func (h *WebhookInteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := Verify(h.server.Logger(), r, h.server.publicKey); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	c := make(chan discord.InteractionResponse)
	go h.server.eventHandlerFunc(discord.GatewayEventTypeInteractionCreate, c, r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(<-c)
	if err != nil {
		h.server.Logger().Errorf("error writing body: %s", err)
	}
}
