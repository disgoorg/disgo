package httpserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
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
	interactionCh    chan channelReader
}

type channelReader struct {
	channel chan discord.InteractionResponse
	reader  io.Reader
}

func (s *serverImpl) Logger() log.Logger {
	return s.config.Logger
}

// PublicKey returns the parsed ed25519.PublicKey
func (s *serverImpl) PublicKey() PublicKey {
	return s.publicKey
}

func (s *serverImpl) Handle(c chan discord.InteractionResponse, payload io.Reader) {
	s.interactionCh <- channelReader{
		channel: c,
		reader:  payload,
	}
}

func (s *serverImpl) listen() {
	s.interactionCh = make(chan channelReader)
	go func() {
		for reader := range s.interactionCh {
			s.eventHandlerFunc(reader.channel, reader.reader)
		}
	}()
}

// Start makes the serverImpl listen on the specified port and handle requests
func (s *serverImpl) Start() error {
	s.config.ServeMux.Handle(s.config.URL, &WebhookInteractionHandler{server: s})
	s.config.HTTPServer.Addr = s.config.Address
	s.config.HTTPServer.Handler = s.config.ServeMux

	s.listen()

	if s.config.CertFile != "" && s.config.KeyFile != "" {
		return s.config.HTTPServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
	}
	return s.config.HTTPServer.ListenAndServe()
}

// Close shuts down the serverImpl
func (s *serverImpl) Close(ctx context.Context) {
	_ = s.config.HTTPServer.Shutdown(ctx)
	close(s.interactionCh)
}

type WebhookInteractionHandler struct {
	server Server
}

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

	responseChannel := make(chan discord.InteractionResponse, 1)
	go h.server.Handle(responseChannel, rqBody)

	timer := time.NewTimer(time.Second * 3)
	defer timer.Stop()
	var (
		body any
		err  error
	)
	select {
	case response := <-responseChannel:
		body, err = response.ToBody()
		h.server.Logger().Info("received response from event handler")
	case <-timer.C:
		h.server.Logger().Warn("interaction timed out")
		http.Error(w, "interaction timed out", http.StatusRequestTimeout)
		return
	}

	if err != nil {
		h.server.Logger().Error("error while converting interaction response to body: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		h.server.Logger().Error("error writing http interaction response error: ", err)
		return
	}

	rsData, _ := ioutil.ReadAll(rsBody)
	h.server.Logger().Trace("response to http interaction. body: ", string(rsData))
}
