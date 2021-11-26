package httpserver

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"io/ioutil"
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

	server := &serverImpl{
		publicKey: hexDecodedKey,
	}

	if config.HTTPServer == nil {
		config.HTTPServer = &http.Server{
			Addr: config.Port,
		}
	}
	server.server = config.HTTPServer

	if config.HTTPServer.Handler == nil {
		if config.ServeMux == nil {
			config.ServeMux = http.NewServeMux()
		}
		config.ServeMux.Handle(config.URL, &WebhookInteractionHandler{server: server})
		config.HTTPServer.Handler = config.ServeMux
	}
	server.config = *config
	return server
}

// serverImpl is used in Bot's webhook server for interactions
type serverImpl struct {
	config    Config
	publicKey ed25519.PublicKey
	server    *http.Server
}

func (s *serverImpl) Logger() log.Logger {
	return s.config.Logger
}

// PublicKey returns the parsed ed25519.PublicKey
func (s *serverImpl) PublicKey() ed25519.PublicKey {
	return s.publicKey
}

// Config returns the Config
func (s *serverImpl) Config() Config {
	return s.config
}

// Start makes the serverImpl listen on the specified port and handle requests
func (s *serverImpl) Start() {
	go func() {
		var err error
		if s.config.CertFile != "" && s.config.KeyFile != "" {
			err = s.server.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil {
			s.Logger().Error("error starting http server: ", err)
		}
	}()
}

// Close shuts down the serverImpl
func (s *serverImpl) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

type WebhookInteractionHandler struct {
	server Server
}

func (h *WebhookInteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := Verify(h.server.Logger(), r, h.server.PublicKey()); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		data, _ := ioutil.ReadAll(r.Body)
		h.server.Logger().Debug("received http interaction with invalid signature. body: ", string(data))
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	body := &bytes.Buffer{}
	data, _ := ioutil.ReadAll(io.TeeReader(r.Body, body))

	h.server.Logger().Debug("received http interaction. body: ", string(data))

	responseChannel := make(chan discord.InteractionResponse, 1)
	h.server.Config().EventHandlerFunc(responseChannel, body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(<-responseChannel)
	if err != nil {
		h.server.Logger().Error("error writing body: ", err)
	}
}
