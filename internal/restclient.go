package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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

// Request makes a new rest request to discords api with the specific route
func (r RestClientImpl) Request(route endpoints.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) error {
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

	rq, err := http.NewRequest(route.Method().String(), route.Route(), reader)
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

	var rawRsBody []byte
	if rsBody == nil {
		rawRsBody = nil
	} else {
		rawRsBody, err = ioutil.ReadAll(rs.Body)
		if err != nil {
			log.Errorf("error reading from response body: %s", err)
			return err
		}
	}

	log.Debugf("code: %d, response: %s", rs.StatusCode, string(rawRsBody))

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if rsBody != nil {
			if err = json.Unmarshal(rawRsBody, rsBody); err != nil {
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
		if err = json.Unmarshal(rawRsBody, &errorRs); err != nil {
			log.Errorf("error unmarshalling error response. code: %d, error: %s", rs.StatusCode, err)
			return err
		}
		return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message_events: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)
	}
}
// GetUser fetches the specific user
func (r RestClientImpl) GetUser(userID api.Snowflake) (user *api.User, err error) {
	return user, r.Request(endpoints.GetUser.Compile(userID), nil, user)
}
// GetMember fetches the specific member
func (r RestClientImpl) GetMember(guildID api.Snowflake, userId api.Snowflake) (member *api.Member, err error) {
	err = r.Request(endpoints.GetMember.Compile(guildID, userId), nil, member)
	return
}
// SendMessage lets you send a message_events to a channel
func (r RestClientImpl) SendMessage(channelID api.Snowflake, message api.Message) *promise.Promise {
	return r.RequestAsync(endpoints.CreateMessage.Compile(channelID), message, &api.Message{})
}
// OpenDMChannel opens a new dm channel a user
func (r RestClientImpl) OpenDMChannel(userId api.Snowflake) *promise.Promise {
	body := struct {RecipientID api.Snowflake `json:"recipient_id"`}{RecipientID: userId}
	return r.RequestAsync(endpoints.PostUsersMeChannels.Compile(), body, &api.DMChannel{})
}
// AddReaction lets you add a reaction to a message_events
func (r RestClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.PutReaction.Compile(channelID, messageID, emoji), nil, nil)
}
// RemoveOwnReaction lets you remove your own reaction from a message_events
func (r RestClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.DeleteOwnReaction.Compile(channelID, messageID, emoji), nil, nil)
}
// RemoveUserReaction lets you remove a specific reaction from a user from a message_events
func (r RestClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) *promise.Promise {
	emoji = strings.Replace(emoji, "#", "%23", -1)
	return r.RequestAsync(endpoints.DeleteUserReaction.Compile(channelID, messageID, emoji, userID), nil, nil)
}

// GetGlobalApplicationCommands gets you all global commands
func (r RestClientImpl) GetGlobalApplicationCommands(applicationID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommands.Compile(applicationID), nil, []api.ApplicationCommand{})
}
// CreateGlobalApplicationGlobalCommand lets you create a new global command
func (r RestClientImpl) CreateGlobalApplicationGlobalCommand(applicationID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.CreateGlobalApplicationCommand.Compile(applicationID), command, api.ApplicationCommand{})
}
// SetGlobalApplicationCommands lets you override all global commands
func (r RestClientImpl) SetGlobalApplicationCommands(applicationID api.Snowflake, commands ...api.ApplicationCommand) *promise.Promise {
	if len(commands) > 100 {
		return promise.Reject(api.ErrTooMuchApplicationCommands)
	}
	return r.RequestAsync(endpoints.SetGlobalApplicationCommands.Compile(applicationID), commands, []api.ApplicationCommand{})
}
// GetGlobalApplicationCommand gets you a specific global global command
func (r RestClientImpl) GetGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommand.Compile(applicationID, commandID), nil, api.ApplicationCommand{})
}
// EditGlobalApplicationCommand lets you edit a specific global command
func (r RestClientImpl) EditGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.EditGlobalApplicationCommand.Compile(applicationID, commandID), command, api.ApplicationCommand{})
}
// DeleteGlobalApplicationCommand lets you delete a specific global command
func (r RestClientImpl) DeleteGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.DeleteGlobalApplicationCommand.Compile(applicationID, commandID), nil, nil)
}

// GetGuildApplicationCommands gets you all guild_events commands
func (r RestClientImpl) GetGuildApplicationCommands(applicationID api.Snowflake, guildID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommands.Compile(applicationID, guildID), nil, []api.ApplicationCommand{})
}
// CreateGuildApplicationGuildCommand lets you create a new guild_events command
func (r RestClientImpl) CreateGuildApplicationGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.CreateGlobalApplicationCommand.Compile(applicationID, guildID), command, api.ApplicationCommand{})
}
// SetGuildApplicationCommands lets you override all guild_events commands
func (r RestClientImpl) SetGuildApplicationCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.ApplicationCommand) *promise.Promise {
	if len(commands) > 100 {
		return promise.Reject(api.ErrTooMuchApplicationCommands)
	}
	return r.RequestAsync(endpoints.SetGlobalApplicationCommands.Compile(applicationID, guildID), commands, []api.ApplicationCommand{})
}
// GetGuildApplicationCommand gets you a specific guild_events command
func (r RestClientImpl) GetGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.GetGlobalApplicationCommand.Compile(applicationID, guildID, commandID), nil, api.ApplicationCommand{})
}
// EditGuildApplicationCommand lets you edit a specific guild_events command
func (r RestClientImpl) EditGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.ApplicationCommand) *promise.Promise {
	return r.RequestAsync(endpoints.EditGlobalApplicationCommand.Compile(applicationID, guildID, commandID), command, api.ApplicationCommand{})
}
// DeleteGuildApplicationCommand lets you delete a specific guild_events command
func (r RestClientImpl) DeleteGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.DeleteGlobalApplicationCommand.Compile(applicationID, guildID, commandID), nil, nil)
}


// Interaction Responses
// SendInteractionResponse used to send the initial response on an interaction
func (r RestClientImpl) SendInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse.Compile(interactionID, interactionToken), interactionResponse, nil)
}
// EditInteractionResponse used to edit the initial response on an interaction
func (r RestClientImpl) EditInteractionResponse(applicationID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) *promise.Promise {
	return r.RequestAsync(endpoints.EditInteractionResponse.Compile(applicationID, interactionToken), interactionResponse, nil)
}
// DeleteInteractionResponse used to delete the initial response on an interaction
func (r RestClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) *promise.Promise {
	return r.RequestAsync(endpoints.DeleteInteractionResponse.Compile(applicationID, interactionToken), nil, nil)
}
// SendFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) SendFollowupMessage(applicationID api.Snowflake, interactionToken string, followupMessage api.FollowupMessage) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken), followupMessage, nil)
}
// EditFollowupMessage used to send the initial response on an interaction
func (r RestClientImpl) EditFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, followupMessage api.InteractionResponse) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken, messageID), followupMessage, nil)
}
// DeleteFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) *promise.Promise {
	return r.RequestAsync(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken, messageID), nil, nil)
}

func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}
