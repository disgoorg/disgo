package internal

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/DisgoOrg/disgo/api"
	"github.com/gorilla/mux"
)

func newWebhookServerImpl(disgo api.Disgo, listenURL string, listenPort int, publicKey string) api.WebhookServer {
	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		disgo.Logger().Errorf("error while decoding hex string: %s", err)
	}
	w := &WebhookServerImpl{
		disgo:      disgo,
		publicKey:  ed25519.PublicKey(hexDecodedKey),
		listenURL:  listenURL,
		listenPort: listenPort,
	}

	w.interactionHandler = &webhookInteractionHandler{disgo: disgo, webhookServer: w}

	return w
}

// WebhookServerImpl is used in Disgo's webhook server for interactions
type WebhookServerImpl struct {
	disgo              api.Disgo
	publicKey          ed25519.PublicKey
	listenURL          string
	listenPort         int
	interactionHandler http.Handler
	router             *mux.Router
}

// Disgo returns the Disgo client
func (w *WebhookServerImpl) Disgo() api.Disgo {
	return w.disgo
}

// PublicKey returns the public key to verify the discord requests with
func (w *WebhookServerImpl) PublicKey() ed25519.PublicKey {
	return w.publicKey
}

// ListenURL returns the URL that the server is listening on
func (w *WebhookServerImpl) ListenURL() string {
	return w.listenURL
}

// Router returns the mux router used on the webhook server
func (w *WebhookServerImpl) Router() *mux.Router {
	return w.router
}

// Start makes the WebhookServerImpl listen on the specified port and handle requests
func (w *WebhookServerImpl) Start() {
	w.router = mux.NewRouter()
	w.router.Handle(w.ListenURL(), w.interactionHandler).Methods("POST")

	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(w.listenPort), w.router); err != nil {
			w.Disgo().Logger().Errorf("error starting webhook server: %s", err)
		}
	}()
}

// Close shuts down the WebhookServerImpl
func (w *WebhookServerImpl) Close() {
	// TODO somehow shutdown the http server gracefully
}

type webhookInteractionHandler struct {
	disgo         api.Disgo
	webhookServer api.WebhookServer
}

func (h *webhookInteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := api.Verify(h.disgo, r, h.webhookServer.PublicKey()); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.disgo.Logger().Errorf("error closing request body in WebhookServer: %s", err)
		}
	}()
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.disgo.Logger().Errorf("error reading body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := make(chan *api.InteractionResponse)
	go h.webhookServer.Disgo().EventManager().Handle(api.WebhookEventInteractionCreate, c, -1, rawBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	err = json.NewEncoder(w).Encode(<-c)
	if err != nil {
		h.disgo.Logger().Errorf("error writing body: %s", err)
	}
}
