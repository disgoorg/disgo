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

// RestClientImpl is the rest client implementation used for HTTP requests to discord
type RestClientImpl struct {
	token string
	Client      *http.Client
}

// Close cleans up the http managers connections
func (r RestClientImpl) Close() {
	r.Client.CloseIdleConnections()
}

// RequestAsync makes a new rest request async to discords api with the specific route
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
func (r RestClientImpl) Request(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) error {
	var reader io.Reader
	var rqJSON []byte
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
		return api.ErrRatelimited

	case http.StatusBadGateway:
		return api.ErrBadGateway

	case http.StatusBadRequest:
		log.Errorf("bad request: %s", string(rqJSON))
		return  api.ErrBadRequest

	case http.StatusUnauthorized:
		return api.ErrUnauthorized

	default:
		var errorRs api.ErrorResponse
		if err = json.Unmarshal(rsBody, &errorRs); err != nil {
			log.Errorf("error unmarshalling error response. code: %d, error: %s", rs.StatusCode, err)
			return err
		}
		return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message_events: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)
	}
}
// GetUser fetches the specific user
func (r RestClientImpl) GetUser(userID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetUser, nil, &api.User{}, userID.String())
}
// GetMember fetches the specific member
func (r RestClientImpl) GetMember(guildID api.Snowflake, userId api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetMember, nil, &api.Member{}, guildID.String(), userId.String())
}
// SendMessage lets you send a message_events to a channel
func (r RestClientImpl) SendMessage(channelID api.Snowflake, message api.Message) *promise.Promise {
	return r.RequestAsync(endpoints.CreateMessage, message, &api.Message{}, channelID.String())
}
// OpenDMChannel opens a new dm channel a user
func (r RestClientImpl) OpenDMChannel(userId api.Snowflake) *promise.Promise {
	body := struct {RecipientID api.Snowflake `json:"recipient_id"`}{RecipientID: userId}
	return r.RequestAsync(endpoints.PostUsersMeChannels, body, &api.DMChannel{})
}
// AddReaction lets you add a reaction to a message_events
func (r RestClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.PutReaction, nil, nil, channelID.String(), messageID.String(), emoji)
}
// RemoveOwnReaction lets you remove your own reaction from a message_events
func (r RestClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.DeleteOwnReaction, nil, nil, channelID.String(), messageID.String(), emoji)
}
// RemoveUserReaction lets you remove a specific reaction from a user from a message_events
func (r RestClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.DeleteUserReaction, nil, nil, channelID.String(), messageID.String(), emoji, userID.String())
}

// GetGlobalApplicationCommands gets you all global commands
func (r RestClientImpl) GetGlobalApplicationCommands(applicationID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommands, nil, []api.ApplicationCommand{}, applicationID.String())
}
// CreateGlobalApplicationGlobalCommand lets you create a new global command
func (r RestClientImpl) CreateGlobalApplicationGlobalCommand(applicationID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.CreateGlobalApplicationCommand, command, api.ApplicationCommand{}, applicationID.String())
}
// SetGlobalApplicationCommands lets you override all global commands
func (r RestClientImpl) SetGlobalApplicationCommands(applicationID api.Snowflake, commands ...api.ApplicationCommand) *promise.Promise {
	if len(commands) > 100 {
		return promise.Reject(api.ErrTooMuchApplicationCommands)
	}
	return r.RequestAsync(endpoints.SetGlobalApplicationCommands, commands, []api.ApplicationCommand{}, applicationID.String())
}
// GetGlobalApplicationCommand gets you a specific global global command
func (r RestClientImpl) GetGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommand, nil, api.ApplicationCommand{}, applicationID.String(), commandID.String())
}
// EditGlobalApplicationCommand lets you edit a specific global command
func (r RestClientImpl) EditGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.EditGlobalApplicationCommand, command, api.ApplicationCommand{}, applicationID.String(), commandID.String())
}
// DeleteGlobalApplicationCommand lets you delete a specific global command
func (r RestClientImpl) DeleteGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.DeleteGlobalApplicationCommand, nil, nil, applicationID.String(), commandID.String())
}

// GetGuildApplicationCommands gets you all guild_events commands
func (r RestClientImpl) GetGuildApplicationCommands(applicationID api.Snowflake, guildID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommands, nil, []api.ApplicationCommand{}, applicationID.String(), guildID.String())
}
// CreateGuildApplicationGuildCommand lets you create a new guild_events command
func (r RestClientImpl) CreateGuildApplicationGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.CreateGlobalApplicationCommand, command, api.ApplicationCommand{}, applicationID.String(), guildID.String())
}
// SetGuildApplicationCommands lets you override all guild_events commands
func (r RestClientImpl) SetGuildApplicationCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.ApplicationCommand) *promise.Promise {
	if len(commands) > 100 {
		return promise.Reject(api.ErrTooMuchApplicationCommands)
	}
	return r.RequestAsync(endpoints.SetGlobalApplicationCommands, commands, []api.ApplicationCommand{}, applicationID.String(), guildID.String())
}
// GetGuildApplicationCommand gets you a specific guild_events command
func (r RestClientImpl) GetGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommand, nil, api.ApplicationCommand{}, applicationID.String(), guildID.String(), commandID.String())
}
// EditGuildApplicationCommand lets you edit a specific guild_events command
func (r RestClientImpl) EditGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.EditGlobalApplicationCommand, command, api.ApplicationCommand{}, applicationID.String(), guildID.String(), commandID.String())
}
// DeleteGuildApplicationCommand lets you delete a specific guild_events command
func (r RestClientImpl) DeleteGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.DeleteGlobalApplicationCommand, nil, nil, applicationID.String(), guildID.String(), commandID.String())
}


// Interaction Responses
// SendInteractionResponse used to send the initial response on an interaction
func (r RestClientImpl) SendInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse, interactionResponse, nil, interactionID.String(), interactionToken)
}
// EditInteractionResponse used to edit the initial response on an interaction
func (r RestClientImpl) EditInteractionResponse(applicationID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) *promise.Promise {
	return r.RequestAsync(endpoints.EditInteractionResponse, interactionResponse, nil, applicationID.String(), interactionToken)
}
// DeleteInteractionResponse used to delete the initial response on an interaction
func (r RestClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) *promise.Promise {
	return r.RequestAsync(endpoints.DeleteInteractionResponse, nil, nil, applicationID.String(), interactionToken)
}
// SendFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) SendFollowupMessage(applicationID api.Snowflake, interactionToken string, followupMessage api.FollowupMessage) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse, followupMessage, nil, applicationID.String(), interactionToken)
}
// EditFollowupMessage used to send the initial response on an interaction
func (r RestClientImpl) EditFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, followupMessage api.InteractionResponse) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse, followupMessage, nil, applicationID.String(), interactionToken, messageID.String())
}
// DeleteFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse, nil, nil, applicationID.String(), interactionToken, messageID.String())
}

func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}
