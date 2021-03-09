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
	"github.com/DiscoOrg/disgo/api/models"
)

// RestClient is the client used for HTTP requests to discord
type RestClientImpl struct {
	DisgoClient api.Disgo
	Client      *http.Client
}

// Disgo returns the Disgo client
func (r RestClientImpl) Disgo() api.Disgo {
	return r.DisgoClient
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
func (r RestClientImpl) RequestAsync(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) *promise.Promise {
	return promise.New(func(resolve func(promise.Any), reject func(error)) {
		err := r.Request(route, rqBody, v, args...)
		if err != nil {
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

	rq.Header.Set("GetUser-Agent", r.UserAgent())
	rq.Header.Set("Authorization", "Bot "+r.Disgo().Token())
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

	log.Infof("code: %d, response: %s", rs.StatusCode, string(rsBody))

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

func (r RestClientImpl) GetUserById(userID models.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetUser, nil, &models.User{}, userID)
}

func (r RestClientImpl) GetMemberById(guildID models.Snowflake, userId models.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetMember, nil, &models.Member{}, guildID, userId)
}

func (r RestClientImpl) SendMessage(channelID models.Snowflake, message models.Message) *promise.Promise {
	return r.RequestAsync(endpoints.PostMessage, message, &models.Message{}, channelID)
}

func (r RestClientImpl) OpenDMChannel(userId models.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.PostUsersMeChannels, models.CreateDMChannel{RecipientID: userId}, &models.DMChannel{})
}

func (r RestClientImpl) AddReaction(channelID models.Snowflake, messageID models.Snowflake, emoji string) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.PutReaction, nil, nil, channelID, messageID, emoji)
}

func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}
