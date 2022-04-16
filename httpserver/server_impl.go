package httpserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/disgo/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

var _ Server = (*serverImpl)(nil)

func New(eventHandlerFunc EventHandlerFunc, opts ...ConfigOpt) Server {
	config := DefaultConfig()
	config.Apply(opts)

	hexDecodedKey, err := hex.DecodeString(config.PublicKey)
	if err != nil {
		config.Logger.Errorf("error while decoding hex string: %s", err)
	}

	return &serverImpl{
		config:           *config,
		publicKey:        hexDecodedKey,
		eventHandlerFunc: eventHandlerFunc,
	}
}

// serverImpl is used in Client's webhook server for interactions
type serverImpl struct {
	config           Config
	publicKey        PublicKey
	eventHandlerFunc EventHandlerFunc
}

func (s *serverImpl) Logger() log.Logger {
	return s.config.Logger
}

// PublicKey returns the parsed ed25519.PublicKey
func (s *serverImpl) PublicKey() PublicKey {
	return s.publicKey
}

func (s *serverImpl) Handle(respondFunc RespondFunc, payload io.Reader) {
	s.eventHandlerFunc(respondFunc, payload)
}

// Start makes the serverImpl listen on the specified port and handle requests
func (s *serverImpl) Start() {
	s.config.ServeMux.Handle(s.config.URL, &WebhookInteractionHandler{server: s})
	s.config.HTTPServer.Addr = s.config.Address
	s.config.HTTPServer.Handler = s.config.ServeMux

	go func() {
		var err error
		if s.config.CertFile != "" && s.config.KeyFile != "" {
			err = s.config.HTTPServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
		} else {
			err = s.config.HTTPServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.Logger().Error("error while running http server: ", err)
		}
	}()
}

// Close shuts down the serverImpl
func (s *serverImpl) Close(ctx context.Context) {
	_ = s.config.HTTPServer.Shutdown(ctx)
}

type WebhookInteractionHandler struct {
	server Server
}

type replyStatus int

const (
	replyStatusWaiting replyStatus = iota
	replyStatusReplied
	replyStatusTimedOut
)

func (h *WebhookInteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := VerifyRequest(h.server.Logger(), r, h.server.PublicKey()); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		data, _ := ioutil.ReadAll(r.Body)
		h.server.Logger().Trace("received http interaction with invalid signature. body: ", string(data))
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	rqBody := &bytes.Buffer{}
	rqData, _ := ioutil.ReadAll(io.TeeReader(r.Body, rqBody))
	h.server.Logger().Trace("received http interaction. body: ", string(rqData))

	// these channels are used to communicate between the http handler and where the interaction is responded to
	responseChannel := make(chan discord.InteractionResponse)
	defer close(responseChannel)
	errorChannel := make(chan error)
	defer close(errorChannel)

	// status of this interaction with a mutex to ensure usage between multiple goroutines
	var (
		status replyStatus
		mu     sync.Mutex
	)

	// send interaction to our handler
	go h.server.Handle(func(response discord.InteractionResponse) error {
		mu.Lock()
		defer mu.Unlock()

		if status == replyStatusTimedOut {
			return discord.ErrInteractionExpired
		}

		if status == replyStatusReplied {
			return discord.ErrInteractionAlreadyReplied
		}

		status = replyStatusReplied
		responseChannel <- response
		// wait if we get any error while processing the response
		return <-errorChannel
	}, rqBody)

	var (
		body any
		err  error
	)

	// wait for the interaction to be responded to or to time out after 3s
	timer := time.NewTimer(time.Second * 3)
	defer timer.Stop()
	select {
	case response := <-responseChannel:
		if body, err = response.ToBody(); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			errorChannel <- err
			return
		}

	case <-timer.C:
		mu.Lock()
		defer mu.Unlock()
		status = replyStatusTimedOut

		h.server.Logger().Debug("interaction timed out")
		http.Error(w, "interaction timed out", http.StatusRequestTimeout)
		return
	}

	rsBody := &bytes.Buffer{}
	multiWriter := io.MultiWriter(w, rsBody)

	if multiPart, ok := body.(*discord.MultipartBuffer); ok {
		w.Header().Set("Content-Type", multiPart.ContentType)
		_, err = io.Copy(multiWriter, multiPart.Buffer)
	} else {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(multiWriter).Encode(body)
	}
	if err != nil {
		errorChannel <- err
		return
	}

	rsData, _ := ioutil.ReadAll(rsBody)
	h.server.Logger().Trace("response to http interaction. body: ", string(rsData))
}
