package internal

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
)

func newWebhookServerImpl(disgo api.Disgo, listenURL string, listenPort int, publicKey string) api.WebhookServer {
	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Errorf("error while decoding hex string: %s", err)
	}
	w := &WebhookServerImpl{
		disgo:      disgo,
		publicKey:  ed25519.PublicKey(hexDecodedKey),
		listenURL:  listenURL,
		listenPort: listenPort,
	}

	w.interactionHandler = &webhookInteractionHandler{w}

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
func (w *WebhookServerImpl) Start() error {
	w.router = mux.NewRouter()
	w.router.Handle(w.ListenURL(), w.interactionHandler).Methods("POST")

	go http.ListenAndServe(":"+strconv.Itoa(w.listenPort), w.router)
	return nil
}

// Close shuts down the WebhookServerImpl
func (w *WebhookServerImpl) Close() {
	// TODO somehow shutdown the http server gracefully
}

type webhookInteractionHandler struct {
	webhookServer api.WebhookServer
}

func (h *webhookInteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := api.Verify(r, h.webhookServer.PublicKey()); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	c := make(chan interface{})
	h.webhookServer.Disgo().EventManager().Handle(api.InteractionCreateWebhookEvent, rawBody, c)

	response := <-c
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorf("error writing body: %s", err)
	}
}
