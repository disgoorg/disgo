package httpserver

import (
	"crypto/ed25519"
	"encoding/hex"
	"net/http"

	"github.com/DisgoOrg/disgo/json"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

var _ Server = (*serverImpl)(nil)

func New(eventHandlerFunc EventHandlerFunc, config *Config) Server {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	config.EventHandlerFunc = eventHandlerFunc

	hexDecodedKey, err := hex.DecodeString(config.PublicKey)
	if err != nil {
		config.Logger.Errorf("error while decoding hex string: %s", err)
	}

	mux := http.NewServeMux()
	server := &serverImpl{
		config:    *config,
		publicKey: hexDecodedKey,
		server: &http.Server{
			Addr:    config.Port,
			Handler: mux,
		},
	}

	mux.Handle(config.URL, &WebhookInteractionHandler{server: server})
	return server
}

// serverImpl is used in Bot's webhook server for interactions
type serverImpl struct {
	config    Config
	publicKey ed25519.PublicKey
	server    *http.Server
}

func (w *serverImpl) Logger() log.Logger {
	return w.config.Logger
}

// Config returns the Config
func (w *serverImpl) Config() Config {
	return w.config
}

// Start makes the serverImpl listen on the specified port and handle requests
func (w *serverImpl) Start() {
	go func() {
		var err error
		if w.config.CertFile != "" && w.config.KeyFile != "" {
			err = w.server.ListenAndServeTLS(w.config.CertFile, w.config.KeyFile)
		} else {
			err = w.server.ListenAndServe()
		}
		if err != nil {
			w.Logger().Errorf("error starting http server: %s", err)
		}
	}()
}

// Close shuts down the serverImpl
func (w *serverImpl) Close() {
	if err := w.server.Close(); err != nil {
		w.Logger().Errorf("error while shutting down http server: %s", err)
	}
}

type WebhookInteractionHandler struct {
	server *serverImpl
}

func (h *WebhookInteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := Verify(h.server.Logger(), r, h.server.publicKey); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	responseChannel := make(chan discord.InteractionResponse, 1)
	h.server.config.EventHandlerFunc(responseChannel, r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(<-responseChannel)
	if err != nil {
		h.server.Logger().Errorf("error writing body: %s", err)
	}
}
