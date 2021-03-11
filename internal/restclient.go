package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chebyrash/promise"
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/endpoints"
)

func newRestClientImpl(token string) api.RestClient {
	return &RestClientImpl{
		token:  token,
		Client: &http.Client{},
	}
}

// RestClient is the client used for HTTP requests to discord
type RestClientImpl struct {
	token string
	Client      *http.Client
}

// Close cleans up the http managers connections
func (r RestClientImpl) Close() {
	r.Client.CloseIdleConnections()
}

// RequestAsync makes a new rest request async to discords api with the specific route
// route the route to make a request to
// rqBody the request body if one should be sent
// v the struct which the response should be unmarshalled in
// args path params to compile the route
func (r RestClientImpl) RequestAsync(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) *promise.Promise {
	return promise.New(func(resolve func(promise.Any), reject func(error)) {
		err := r.Request(route, rqBody, v, args...)
		if err != nil {
			log.Errorf("received error on route: %s. error: %s", route.Compile(args...), err)
			reject(err)
			return
		}
		resolve(v)
	})
}

// Request makes a new rest request to discords api with the specific route
// route the route to make a request to
// rqBody the request body if one should be sent
// v the struct which the response should be unmarshalled in
// args path params to compile the route
func (r RestClientImpl) Request(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) error {
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

	rq, err := http.NewRequest(route.Method().String(), route.Compile(args...), reader)
	if err != nil {
		return err
	}

	rq.Header.Set("GetUser-Agent", r.UserAgent())
	rq.Header.Set("Authorization", "Bot "+r.token)
	rq.Header.Set("content-type", "application/json")

	rs, err := r.Client.Do(rq)
	if err != nil {
		return err
	}

	defer func() {
		err := rs.Body.Close()
		if err != nil {
			log.Error("error closing response body", err.Error())
		}
	}()

	var rsBody []byte
	if v == nil {
		rsBody = nil
	} else {
		rsBody, err = ioutil.ReadAll(rs.Body)
		if err != nil {
			log.Errorf("error reading from response body: %s", err)
			return err
		}
	}

	log.Debugf("code: %d, response: %s", rs.StatusCode, string(rsBody))

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if v != nil {
			if err = json.Unmarshal(rsBody, v); err != nil {
				log.Errorf("error unmarshalling response. error: %s", err)
				return err
			}
		}
		return nil

	case http.StatusTooManyRequests:
		limit := rs.Header.Get("X-RateLimit-Limit")
		remaining := rs.Header.Get("X-RateLimit-Limit")
		reset := rs.Header.Get("X-RateLimit-Limit")
		bucket := rs.Header.Get("X-RateLimit-Limit")
		log.Errorf("too many requests. limit: %s, remaining: %s, reset: %s,bucket: %s", limit, remaining, reset, bucket)
		return api.RatelimitedError

	case http.StatusBadGateway:
		return api.BadGatewayError

	case http.StatusUnauthorized:
		return api.UnauthorizedError

	default:
		var errorRs api.ErrorResponse
		if err = json.Unmarshal(rsBody, &errorRs); err != nil {
			log.Errorf("error unmarshalling error response. error: %s", err)
			return err
		}
		return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)
	}
}

func (r RestClientImpl) GetUserById(userID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetUser, nil, &api.User{}, userID.String())
}

func (r RestClientImpl) GetMemberById(guildID api.Snowflake, userId api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetMember, nil, &api.Member{}, guildID.String(), userId.String())
}

func (r RestClientImpl) SendMessage(channelID api.Snowflake, message api.Message) *promise.Promise {
	return r.RequestAsync(endpoints.CreateMessage, message, &api.Message{}, channelID.String())
}

func (r RestClientImpl) OpenDMChannel(userId api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.PostUsersMeChannels, api.CreateDMChannel{RecipientID: userId}, &api.DMChannel{})
}

func (r RestClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.PutReaction, nil, nil, channelID.String(), messageID.String(), emoji)
}

func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}
