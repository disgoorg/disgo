package httpserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

var _ Server = (*serverImpl)(nil)

// New creates a new Server with the given publicKey eventHandlerFunc and ConfigOpt(s)
func New(publicKey string, eventHandlerFunc EventHandlerFunc, opts ...ConfigOpt) Server {
	config := DefaultConfig()
	config.Apply(opts)

	hexDecodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		config.Logger.Errorf("error while decoding hex string: %s", err)
	}

	return &serverImpl{
		config:           *config,
		publicKey:        hexDecodedKey,
		eventHandlerFunc: eventHandlerFunc,
	}
}

type serverImpl struct {
	config           Config
	publicKey        PublicKey
	eventHandlerFunc EventHandlerFunc
}

func (s *serverImpl) Logger() log.Logger {
	return s.config.Logger
}

func (s *serverImpl) PublicKey() PublicKey {
	return s.publicKey
}

func (s *serverImpl) Handle(respondFunc RespondFunc, event gateway.EventInteractionCreate) {
	s.eventHandlerFunc(respondFunc, event)
}

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

func (s *serverImpl) Close(ctx context.Context) {
	_ = s.config.HTTPServer.Shutdown(ctx)
}

// WebhookInteractionHandler implements the http.Handler interface and is used to handle interactions from Discord.
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
	if ok := VerifyRequest(r, h.server.PublicKey()); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		data, _ := io.ReadAll(r.Body)
		h.server.Logger().Trace("received http interaction with invalid signature. body: ", string(data))
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	buff := new(bytes.Buffer)
	rqData, _ := io.ReadAll(io.TeeReader(r.Body, buff))
	h.server.Logger().Trace("received http interaction. body: ", string(rqData))

	var v gateway.EventInteractionCreate
	if err := json.NewDecoder(buff).Decode(&v); err != nil {
		h.server.Logger().Error("error while decoding interaction: ", err)
		return
	}

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
	}, v)

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

	rsData, _ := io.ReadAll(rsBody)
	h.server.Logger().Trace("response to http interaction. body: ", string(rsData))
}
