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

func newRestClientImpl(disgo api.Disgo, token string) api.RestClient {
	return &RestClientImpl{
		disgo:  disgo,
		Client: &http.Client{},
		token:  token,
	}
}

// RestClientImpl is the rest client implementation used for HTTP requests to discord
type RestClientImpl struct {
	disgo  api.Disgo
	Client *http.Client
	token  string
}

// Close cleans up the http managers connections
func (r RestClientImpl) Disgo() api.Disgo {
	return r.disgo
}

// Close cleans up the http managers connections
func (r RestClientImpl) Close() {
	r.Client.CloseIdleConnections()
}

// UserAgent returns the user agent for this rst client
func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
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

	rq.Header.Set("User-Agent", r.UserAgent())
	rq.Header.Set("Authorization", "Bot "+r.token)
	rq.Header.Set("Content-Type", "application/json")

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
		return api.ErrBadRequest

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

// SendMessage lets you send a message_events to a channel
func (r RestClientImpl) SendMessage(channelID api.Snowflake, message api.Message) (rMessage *api.Message, err error) {
	err = r.Request(endpoints.CreateMessage.Compile(channelID), message, &rMessage)
	if rMessage != nil {
		//r.Disgo().Cache().CacheMessage(rMessage)
	}
	return
}

// OpenDMChannel opens a new dm channel a user
func (r RestClientImpl) OpenDMChannel(userId api.Snowflake) (channel *api.DMChannel, err error) {
	body := struct {
		RecipientID api.Snowflake `json:"recipient_id"`
	}{
		RecipientID: userId,
	}
	err = r.Request(endpoints.CreateDMChannel.Compile(), body, &channel)
	if channel != nil {
		channel.Disgo = r.Disgo()
		r.Disgo().Cache().CacheDMChannel(channel)
	}
	return
}

// UpdateSelfNick updates the bots nickname in a guild
func (r RestClientImpl) UpdateSelfNick(guildID api.Snowflake, nick *string) (newNick *string, err error) {
	var updateNick *api.UpdateSelfNick
	err = r.Request(endpoints.UpdateSelfNick.Compile(guildID), api.UpdateSelfNick{Nick: nick}, &updateNick)
	if updateNick != nil {
		r.Disgo().Cache().Member(guildID, r.Disgo().ApplicationID()).Nick = updateNick.Nick
		newNick = updateNick.Nick
	}
	return
}

// GetUser fetches the specific user
func (r RestClientImpl) GetUser(userID api.Snowflake) (user *api.User, err error) {
	err = r.Request(endpoints.GetUser.Compile(userID), nil, &user)
	if user != nil {
		r.Disgo().Cache().CacheUser(user)
	}
	return
}

// GetMember fetches the specific member
func (r RestClientImpl) GetMember(guildID api.Snowflake, userID api.Snowflake) (member *api.Member, err error) {
	err = r.Request(endpoints.GetMember.Compile(guildID, userID), nil, &member)
	if member != nil {
		r.Disgo().Cache().CacheMember(member)
	}
	return
}

// GetMembers fetches all members for a guild
func (r RestClientImpl) GetMembers(guildID api.Snowflake) (members []*api.Member, err error) {
	err = r.Request(endpoints.GetMembers.Compile(guildID), nil, &members)
	if members != nil {
		for _, member := range members {
			r.Disgo().Cache().CacheMember(member)
		}
	}
	return
}

// AddMember adds a member to the guild with the oauth2 access token. requires api.PermissionCreateInstantInvite
func (r RestClientImpl) AddMember(guildID api.Snowflake, userID api.Snowflake, addGuildMemberData api.AddGuildMemberData) (member *api.Member, err error) {
	err = r.Request(endpoints.AddMember.Compile(guildID, userID), addGuildMemberData, &member)
	if member != nil {
		r.Disgo().Cache().CacheMember(member)
	}
	return
}

// KickMember kicks a member from the guild. requires api.PermissionKickMembers
func (r RestClientImpl) KickMember(guildID api.Snowflake, userID api.Snowflake, reason *string) (err error) {
	var route endpoints.CompiledAPIRoute
	if reason == nil {
		route = endpoints.RemoveMember.Compile(guildID, userID)
	} else {
		route = endpoints.RemoveMemberReason.Compile(guildID, userID, *reason)
	}

	err = r.Request(route, nil, nil)
	if err == nil {
		r.Disgo().Cache().UncacheMember(guildID, userID)
	}
	return
}

// UpdateMember updates a member
func (r RestClientImpl) UpdateMember(guildID api.Snowflake, userID api.Snowflake, updateGuildMemberData api.UpdateGuildMemberData) (member *api.Member, err error) {
	err = r.Request(endpoints.UpdateMember.Compile(guildID, userID), updateGuildMemberData, &member)
	if member != nil {
		r.Disgo().Cache().CacheMember(member)
	}
	return
}

// MoveMember moves/kicks the member to/from a voice channel
func (r RestClientImpl) MoveMember(guildID api.Snowflake, userID api.Snowflake, channelID *api.Snowflake) (member *api.Member, err error) {
	err = r.Request(endpoints.UpdateMember.Compile(guildID, userID), api.MoveGuildMemberData{ChannelID: channelID}, &member)
	if member != nil {
		r.Disgo().Cache().CacheMember(member)
	}
	return
}

// AddMemberRole adds a role to a member
func (r RestClientImpl) AddMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
	err = r.Request(endpoints.AddMemberRole.Compile(guildID, userID, roleID), nil, nil)
	if err == nil {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			member.Roles = append(member.Roles, roleID)
		}
	}
	return
}

// AddMemberRole removes a role from a member
func (r RestClientImpl) RemoveMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
	err = r.Request(endpoints.RemoveMemberRole.Compile(guildID, userID, roleID), nil, nil)
	if err == nil {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			for i, id := range member.Roles {
				if id == roleID {
					member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
					break
				}
			}
		}
	}
	return
}

// GetRoles fetches all roles from a guild
func (r RestClientImpl) GetRoles(guildID api.Snowflake) (roles []*api.Role, err error) {
	err = r.Request(endpoints.GetRoles.Compile(guildID), nil, &roles)
	if roles != nil {
		for _, role := range roles {
			role.Disgo = r.disgo
			role.GuildID = guildID
			r.disgo.Cache().CacheRole(role)
		}
	}
	return
}

// CreateRole creates a new role for a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) CreateRole(guildID api.Snowflake, role api.UpdateRole) (newRole *api.Role, err error) {
	err = r.Request(endpoints.CreateRole.Compile(guildID), role, &newRole)
	if newRole != nil {
		newRole.Disgo = r.disgo
		newRole.GuildID = guildID
		r.disgo.Cache().CacheRole(newRole)
	}
	return
}

// UpdateRole updates a role from a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) UpdateRole(guildID api.Snowflake, roleID api.Snowflake, role api.UpdateRole) (newRole *api.Role, err error) {
	err = r.Request(endpoints.UpdateRole.Compile(guildID, roleID), role, &newRole)
	if newRole != nil {
		newRole.Disgo = r.disgo
		newRole.GuildID = guildID
		r.disgo.Cache().CacheRole(newRole)
	}
	return
}

// UpdateRolePositions updates the position of a role from a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) UpdateRolePositions(guildID api.Snowflake, roleUpdates ...api.UpdateRolePosition) (roles []*api.Role, err error) {
	err = r.Request(endpoints.GetRoles.Compile(guildID), roleUpdates, &roles)
	if roles != nil {
		for _, role := range roles {
			role.Disgo = r.disgo
			role.GuildID = guildID
			r.disgo.Cache().CacheRole(role)
		}
	}
	return
}

// DeleteRole deletes a role from a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) DeleteRole(guildID api.Snowflake, roleID api.Snowflake) (err error) {
	err = r.Request(endpoints.UpdateRole.Compile(guildID, roleID), nil, nil)
	if err == nil {
		r.disgo.Cache().UncacheRole(guildID, roleID)
	}
	return
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

// GetGlobalCommands gets you all global commands
func (r RestClientImpl) GetGlobalCommands(applicationID api.Snowflake) (commands []*api.Command, err error) {
	return commands, r.Request(endpoints.GetGlobalApplicationCommands.Compile(applicationID), nil, &commands)
}

// CreateGlobalGlobalCommand lets you create a new global command
func (r RestClientImpl) CreateGlobalCommand(applicationID api.Snowflake, command api.Command) (rCommand *api.Command, err error) {
	return rCommand, r.Request(endpoints.CreateGlobalApplicationCommand.Compile(applicationID), command, &rCommand)
}

// SetGlobalCommands lets you override all global commands
func (r RestClientImpl) SetGlobalCommands(applicationID api.Snowflake, commands ...api.Command) (rCommands []*api.Command, err error) {
	if len(commands) > 100 {
		err = api.ErrTooMuchApplicationCommands
		return
	}
	return rCommands, r.Request(endpoints.SetGlobalApplicationCommands.Compile(applicationID), api.Commands(commands), &rCommands)
}

// GetGlobalCommand gets you a specific global global command
func (r RestClientImpl) GetGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (rCommand *api.Command, err error) {
	return rCommand, r.Request(endpoints.GetGlobalApplicationCommand.Compile(applicationID, commandID), nil, &rCommand)
}

// EditGlobalCommand lets you edit a specific global command
func (r RestClientImpl) EditGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.Command) (rCommand *api.Command, err error) {
	return rCommand, r.Request(endpoints.EditGlobalApplicationCommand.Compile(applicationID, commandID), command, &rCommand)
}

// DeleteGlobalCommand lets you delete a specific global command
func (r RestClientImpl) DeleteGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) error {
	return r.Request(endpoints.DeleteGlobalApplicationCommand.Compile(applicationID, commandID), nil, nil)
}

// GetGuildCommands gets you all guild_events commands
func (r RestClientImpl) GetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake) (commands []*api.Command, err error) {
	return commands, r.Request(endpoints.GetGuildApplicationCommands.Compile(applicationID, guildID), nil, &commands)
}

// CreateGuildGuildCommand lets you create a new guild_events command
func (r RestClientImpl) CreateGuildGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.Command) (rCommand *api.Command, err error) {
	return rCommand, r.Request(endpoints.CreateGuildApplicationCommand.Compile(applicationID, guildID), command, &rCommand)
}

// SetGuildCommands lets you override all guild_events commands
func (r RestClientImpl) SetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.Command) (rCommands []*api.Command, err error) {
	if len(commands) > 100 {
		err = api.ErrTooMuchApplicationCommands
		return
	}
	return rCommands, r.Request(endpoints.SetGuildApplicationCommands.Compile(applicationID, guildID), commands, &rCommands)
}

// GetGuildCommand gets you a specific guild_events command
func (r RestClientImpl) GetGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (rCommand *api.Command, err error) {
	return rCommand, r.Request(endpoints.GetGuildApplicationCommand.Compile(applicationID, guildID, commandID), nil, &rCommand)
}

// EditGuildCommand lets you edit a specific guild_events command
func (r RestClientImpl) EditGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.Command) (rCommand *api.Command, err error) {
	return rCommand, r.Request(endpoints.EditGuildApplicationCommand.Compile(applicationID, guildID, commandID), command, &rCommand)
}

// DeleteGuildCommand lets you delete a specific guild_events command
func (r RestClientImpl) DeleteGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) error {
	return r.Request(endpoints.DeleteGuildApplicationCommand.Compile(applicationID, guildID, commandID), nil, nil)
}

// SendInteractionResponse used to send the initial response on an interaction
func (r RestClientImpl) SendInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) error {
	return r.Request(endpoints.CreateInteractionResponse.Compile(interactionID, interactionToken), interactionResponse, nil)
}

// EditInteractionResponse used to edit the initial response on an interaction
func (r RestClientImpl) EditInteractionResponse(applicationID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) (message *api.Message, err error) {
	return message, r.Request(endpoints.EditInteractionResponse.Compile(applicationID, interactionToken), interactionResponse, &message)
}

// DeleteInteractionResponse used to delete the initial response on an interaction
func (r RestClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) error {
	return r.Request(endpoints.DeleteInteractionResponse.Compile(applicationID, interactionToken), nil, nil)
}

// SendFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) SendFollowupMessage(applicationID api.Snowflake, interactionToken string, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken), followupMessage, &message)
}

// EditFollowupMessage used to send the initial response on an interaction
func (r RestClientImpl) EditFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken, messageID), followupMessage, &message)
}

// DeleteFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) error {
	return r.Request(endpoints.CreateInteractionResponse.Compile(applicationID, interactionToken, messageID), nil, nil)
}

func normalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}
