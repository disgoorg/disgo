package disgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/endpoints"
)

// RestClient is a manager for all of disgo's HTTP requests
type RestClient interface {
	Disgo() Disgo
	Close()
	userAgent() string
	Request(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) error
}

// RestClient is the client used for HTTP requests to discord
type RestClientImpl struct {
	disgo  Disgo
	client *http.Client
}

// Disgo returns the Disgo client
func (r RestClientImpl) Disgo() Disgo {
	return r.disgo
}

// Close cleans up the http managers connections
func (r RestClientImpl) Close() {
	r.client.CloseIdleConnections()
}

// Request makes a new rest request to discords api with the specific route
// route the route to make a request to
// rqBody the request body if one should be sent
// v the struct which the response should be unmarshalled in
// args path params to compile the route
func (r RestClientImpl) Request(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) error {

	var reader io.Reader
	if rqBody != nil {
		rqJSON, err := json.Marshal(rqBody)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer(rqJSON)
	} else {
		reader = nil
	}

	rq, err := http.NewRequest(route.Method.String(), route.Compile(args...), reader)
	if err != nil {
		return err
	}

	rq.Header.Set("User-Agent", r.userAgent())
	rq.Header.Set("Authorization", "Bot "+r.Disgo().Token())
	rq.Header.Set("Content-Type", "application/json")

	rs, err := r.client.Do(rq)
	if err != nil {
		return err
	}

	defer func() {
		err := rs.Body.Close()
		if err != nil {
			log.Error("error closing response body", err.Error())
		}
	}()

	rsBody, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return err
	}

	switch rs.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNoContent:
		if err = json.Unmarshal(rsBody, v); err != nil {
			return err
		}

	case http.StatusTooManyRequests:
		limit := rs.Header.Get("X-RateLimit-Limit")
		remaining := rs.Header.Get("X-RateLimit-Limit")
		reset := rs.Header.Get("X-RateLimit-Limit")
		bucket := rs.Header.Get("X-RateLimit-Limit")
		return fmt.Errorf("too many requests. limit: %s, remaining: %s, reset: %s,bucket: %s", limit, remaining, reset, bucket)

	case http.StatusBadGateway:
		return errors.New("bad gateway could not reach discord")

	case http.StatusUnauthorized:
		return errors.New("the provided token is invalid")

	default:
		var errorRs ErrorResponse
		if err = json.Unmarshal(rsBody, &errorRs); err != nil {
			return err
		}
		return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)
	}

	return nil
}

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}

func (r RestClientImpl) userAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}
