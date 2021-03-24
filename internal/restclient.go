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
	return member, r.Request(endpoints.GetMember.Compile(guildID, userId), nil, member)
}
// SendMessage lets you send a message_events to a channel
func (r RestClientImpl) SendMessage(channelID api.Snowflake, message api.Message) (rMessage *api.Message, err error) {
	return rMessage, r.Request(endpoints.CreateMessage.Compile(channelID), message, rMessage)
}
// OpenDMChannel opens a new dm channel a user
func (r RestClientImpl) OpenDMChannel(userId api.Snowflake) (channel *api.DMChannel, err error) {
	body := struct {RecipientID api.Snowflake `json:"recipient_id"`}{RecipientID: userId}
	return channel, r.Request(endpoints.PostUsersMeChannels.Compile(), body, channel)
}
// AddReaction lets you add a reaction to a message_events
func (r RestClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	return r.Request(endpoints.PutReaction.Compile(channelID, messageID, normalizeEmoji(emoji)), nil, nil)
}
// RemoveOwnReaction lets you remove your own reaction from a message_events
func (r RestClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	return r.Request(endpoints.DeleteOwnReaction.Compile(channelID, messageID, normalizeEmoji(emoji)), nil, nil)
}
// RemoveUserReaction lets you remove a specific reaction from a user from a message_events
func (r RestClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) error {
	return r.Request(endpoints.DeleteUserReaction.Compile(channelID, messageID, normalizeEmoji(emoji), userID), nil, nil)
}

// GetGlobalApplicationCommands gets you all global commands
func (r RestClientImpl) GetGlobalApplicationCommands(applicationID api.Snowflake) (commands []*api.ApplicationCommand, err error) {
	return commands, r.Request(endpoints.GetGlobalApplicationCommands.Compile(applicationID), nil, commands)
}
// CreateGlobalApplicationGlobalCommand lets you create a new global command
func (r RestClientImpl) CreateGlobalApplicationGlobalCommand(applicationID api.Snowflake, command api.ApplicationCommand) (rCommand *api.ApplicationCommand, err error) {
	return rCommand, r.Request(endpoints.CreateGlobalApplicationCommand.Compile(applicationID), command, rCommand)
}
// SetGlobalApplicationCommands lets you override all global commands
func (r RestClientImpl) SetGlobalApplicationCommands(applicationID api.Snowflake, commands ...api.ApplicationCommand) (rCommands []*api.ApplicationCommand, err error) {
	if len(commands) > 100 {
		err = api.ErrTooMuchApplicationCommands
		return
	}
	return rCommands, r.Request(endpoints.SetGlobalApplicationCommands.Compile(applicationID), commands, rCommands)
}
// GetGlobalApplicationCommand gets you a specific global global command
func (r RestClientImpl) GetGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake) (rCommand *api.ApplicationCommand, err error) {
	return rCommand, r.Request(endpoints.GetGlobalApplicationCommand.Compile(applicationID, commandID), nil, rCommand)
}
// EditGlobalApplicationCommand lets you edit a specific global command
func (r RestClientImpl) EditGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.ApplicationCommand) (rCommand *api.ApplicationCommand, err error) {
	return rCommand, r.Request(endpoints.EditGlobalApplicationCommand.Compile(applicationID, commandID), command, rCommand)
}
// DeleteGlobalApplicationCommand lets you delete a specific global command
func (r RestClientImpl) DeleteGlobalApplicationCommand(applicationID api.Snowflake, commandID api.Snowflake) error {
	return r.Request(endpoints.DeleteGlobalApplicationCommand.Compile(applicationID, commandID), nil, nil)
}

// GetGuildApplicationCommands gets you all guild_events commands
func (r RestClientImpl) GetGuildApplicationCommands(applicationID api.Snowflake, guildID api.Snowflake) (commands []*api.ApplicationCommand, err error) {
	return commands, r.Request(endpoints.GetGlobalApplicationCommands.Compile(applicationID, guildID), nil, commands)
}
// CreateGuildApplicationGuildCommand lets you create a new guild_events command
func (r RestClientImpl) CreateGuildApplicationGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.ApplicationCommand) (rCommand *api.ApplicationCommand, err error) {
	return rCommand, r.Request(endpoints.CreateGlobalApplicationCommand.Compile(applicationID, guildID), command, rCommand)
}
// SetGuildApplicationCommands lets you override all guild_events commands
func (r RestClientImpl) SetGuildApplicationCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.ApplicationCommand) (rCommands []*api.ApplicationCommand, err error) {
	if len(commands) > 100 {
		err = api.ErrTooMuchApplicationCommands
		return
	}
	return rCommands, r.Request(endpoints.SetGlobalApplicationCommands.Compile(applicationID, guildID), commands, rCommands)
}
// GetGuildApplicationCommand gets you a specific guild_events command
func (r RestClientImpl) GetGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (rCommand *api.ApplicationCommand, err error) {
	return rCommand, r.Request(endpoints.GetGlobalApplicationCommand.Compile(applicationID, guildID, commandID), nil, rCommand)
}
// EditGuildApplicationCommand lets you edit a specific guild_events command
func (r RestClientImpl) EditGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.ApplicationCommand) (rCommand *api.ApplicationCommand, err error) {
	return rCommand, r.Request(endpoints.EditGlobalApplicationCommand.Compile(applicationID, guildID, commandID), command, rCommand)
}
// DeleteGuildApplicationCommand lets you delete a specific guild_events command
func (r RestClientImpl) DeleteGuildApplicationCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) error {
	return r.Request(endpoints.DeleteGlobalApplicationCommand.Compile(applicationID, guildID, commandID), nil, nil)
}

// SendInteractionResponse used to send the initial response on an interaction
func (r RestClientImpl) SendInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) error {
	return r.Request(endpoints.CreateInteractionResponse.Compile(interactionID, interactionToken), interactionResponse, nil)
}
// EditInteractionResponse used to edit the initial response on an interaction
func (r RestClientImpl) EditInteractionResponse(applicationID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) (message *api.Message, err error) {
	return message, r.Request(endpoints.EditInteractionResponse.Compile(applicationID, interactionToken), interactionResponse, message)
}
// DeleteInteractionResponse used to delete the initial response on an interaction
func (r RestClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) error {
	return r.Request(endpoints.DeleteInteractionResponse.Compile(applicationID, interactionToken), nil, nil)
}

// SendFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) SendFollowupMessage(applicationID api.Snowflake, interactionToken string, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken), followupMessage, message)
}
// EditFollowupMessage used to send the initial response on an interaction
func (r RestClientImpl) EditFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken, messageID), followupMessage, message)
}
// DeleteFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) error {
	return r.Request(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken, messageID), nil, nil)
}

func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}

func normalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}
