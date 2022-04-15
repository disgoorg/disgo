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
		eventHandlerFunc: eventHandlerFunc,
		publicKey:        hexDecodedKey,
	}
}

// serverImpl is used in Client's webhook server for interactions
type serverImpl struct {
	config           Config
	eventHandlerFunc EventHandlerFunc
	publicKey        PublicKey
	interactionCh    chan interaction
}

type interaction struct {
	respondFunc RespondFunc
	reader      io.Reader
}

func (s *serverImpl) Logger() log.Logger {
	return s.config.Logger
}

// PublicKey returns the parsed ed25519.PublicKey
func (s *serverImpl) PublicKey() PublicKey {
	return s.publicKey
}

func (s *serverImpl) Handle(respondFunc RespondFunc, payload io.Reader) {
	s.interactionCh <- interaction{
		respondFunc: respondFunc,
		reader:      payload,
	}
}

func (s *serverImpl) listen() {
	s.interactionCh = make(chan interaction)
	for i := range s.interactionCh {
		s.eventHandlerFunc(i.respondFunc, i.reader)
	}
}

// Start makes the serverImpl listen on the specified port and handle requests
func (s *serverImpl) Start() {
	s.config.ServeMux.Handle(s.config.URL, &WebhookInteractionHandler{server: s})
	s.config.HTTPServer.Addr = s.config.Address
	s.config.HTTPServer.Handler = s.config.ServeMux

	go s.listen()

	go func() {
		var err error
		if s.config.CertFile != "" && s.config.KeyFile != "" {
			err = s.config.HTTPServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
		} else {
			err = s.config.HTTPServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.config.Logger.Error("error while running http server: ", err)
		}
	}()
}

// Close shuts down the serverImpl
func (s *serverImpl) Close(ctx context.Context) {
	_ = s.config.HTTPServer.Shutdown(ctx)
	close(s.interactionCh)
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

	responseChannel := make(chan discord.InteractionResponse)
	defer close(responseChannel)
	errorChannel := make(chan error)
	defer close(errorChannel)
	var (
		status replyStatus
		mu     sync.Mutex
	)

	respondFunc := func(response discord.InteractionResponse) error {
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
	}
	go h.server.Handle(respondFunc, rqBody)

	timer := time.NewTimer(time.Second * 3)
	defer timer.Stop()
	var (
		body any
		err  error
	)
	select {
	case response := <-responseChannel:
		body, err = response.ToBody()
	case <-timer.C:
		mu.Lock()
		defer mu.Unlock()
		status = replyStatusTimedOut

		h.server.Logger().Debug("interaction timed out")
		http.Error(w, "interaction timed out", http.StatusRequestTimeout)
		return
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		errorChannel <- err
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
