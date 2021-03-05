package src

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/src/endpoints"
)

type RestClient struct {
	Client *http.Client
	Token  string
}

// Request makes a new rest request to discords api with the specific route
// route the route to make a request to
// rqBody the request body if one should be sent
// v the struct which the response should be unmarshalled in
// args path params to compile the route
func (c RestClient) Request(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) error {

	rqJson, err := json.Marshal(rqBody)
	if err != nil {
		return err
	}

	rq, err := http.NewRequest(route.Method.String(), route.Compile(args...), bytes.NewBuffer(rqJson))
	if err != nil {
		return err
	}

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

	err = json.Unmarshal(rsBody, v)
	if err != nil {
		return err
	}

	return nil
}
