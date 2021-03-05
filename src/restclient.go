package src

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/src/endpoints"
)

type RestClient struct {
	Client    *http.Client
	Token     string
	UserAgent string
}

// Request makes a new rest request to discords api with the specific route
// route the route to make a request to
// rqBody the request body if one should be sent
// v the struct which the response should be unmarshalled in
// args path params to compile the route
func (c RestClient) Request(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) error {

	var reader io.Reader
	if rqBody != nil {
		rqJson, err := json.Marshal(rqBody)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer(rqJson)
	} else {
		reader = nil
	}

	rq, err := http.NewRequest(route.Method.String(), route.Compile(args...), reader)
	if err != nil {
		return err
	}

	rq.Header.Set("User-Agent", c.UserAgent)
	rq.Header.Set("Authorization", "Bot "+c.Token)
	rq.Header.Set("Content-Type", "application/json")

	rs, err := c.Client.Do(rq)
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

	log.Info(string(rsBody))

	err = json.Unmarshal(rsBody, v)
	if err != nil {
		return err
	}

	return nil
}
