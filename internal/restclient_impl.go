package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DisgoOrg/disgo/api/events"
	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
)

func newRestClientImpl(disgo api.Disgo) api.RestClient {
	return &RestClientImpl{
		disgo:  disgo,
		Client: &http.Client{},
	}
}

// RestClientImpl is the rest client implementation used for HTTP requests to discord
type RestClientImpl struct {
	disgo  api.Disgo
	Client *http.Client
	token  string
}

// Disgo returns the api.Disgo instance
func (r RestClientImpl) Disgo() api.Disgo {
	return r.disgo
}

// Close cleans up the http managers connections
func (r RestClientImpl) Close() {
	r.Client.CloseIdleConnections()
}

// UserAgent returns the user agent for this api.RestClient
func (r RestClientImpl) UserAgent() string {
	return "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)"
}

// Request makes a new rest request to discords api with the specific endpoints.APIRoute
func (r RestClientImpl) Request(route endpoints.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) error {
	var reader io.Reader
	var rqJSON []byte
	if rqBody != nil {
		rqJSON, err := json.Marshal(rqBody)
		if err != nil {
			return err
		}
		log.Debugf("request json: \"%s\"", string(rqJSON))
		reader = bytes.NewBuffer(rqJSON)
	} else {
		reader = nil
	}

	rq, err := http.NewRequest(route.Method().String(), route.Route(), reader)
	if err != nil {
		return err
	}

	rq.Header.Set("User-Agent", r.UserAgent())
	rq.Header.Set("Authorization", string("Bot "+r.disgo.Token()))
	rq.Header.Set("Content-Type", "application/json")

	rs, err := r.Client.Do(rq)
	if err != nil {
		return err
	}

	defer func() {
		err = rs.Body.Close()
		if err != nil {
			log.Error("error closing response body", err.Error())
		}
	}()

	rawRsBody, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		log.Errorf("error reading from response body: %s", err)
		return err
	}

	log.Debugf("code: %d, response: %s", rs.StatusCode, string(rawRsBody))

	r.Disgo().EventManager().Dispatch(events.HttpRequestEvent{
		GenericEvent: events.NewEvent(r.Disgo(), 0),
		Request:      rq,
		Response:     rs,
	})

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
		log.Error(api.ErrBadGateway)
		return api.ErrBadGateway

	case http.StatusBadRequest:
		log.Errorf("bad request request: \"%s\", response: \"%s\"", string(rqJSON), string(rawRsBody))
		return api.ErrBadRequest

	case http.StatusUnauthorized:
		log.Error(api.ErrUnauthorized)
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

// SendMessage lets you send a api.Message to a api.MessageChannel
func (r RestClientImpl) SendMessage(channelID api.Snowflake, message api.MessageCreate) (msg *api.Message, err error) {
	err = r.Request(endpoints.CreateMessage.Compile(channelID), message, &msg)
	if err == nil {
		msg = r.Disgo().EntityBuilder().CreateMessage(msg, api.CacheStrategyNoWs)
	}
	return
}

// EditMessage lets you edit a api.Message
func (r RestClientImpl) EditMessage(channelID api.Snowflake, messageID api.Snowflake, message api.MessageUpdate) (msg *api.Message, err error) {
	err = r.Request(endpoints.UpdateMessage.Compile(channelID, messageID), message, &msg)
	if err == nil {
		msg = r.Disgo().EntityBuilder().CreateMessage(msg, api.CacheStrategyNoWs)
	}
	return
}

// DeleteMessage lets you delete a api.Message
func (r RestClientImpl) DeleteMessage(channelID api.Snowflake, messageID api.Snowflake) (err error) {
	err = r.Request(endpoints.DeleteMessage.Compile(channelID, messageID), nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMessage(channelID, messageID)
	}
	return
}

// BulkDeleteMessages lets you bulk delete api.Message(s)
func (r RestClientImpl) BulkDeleteMessages(channelID api.Snowflake, messageIDs ...api.Snowflake) (err error) {
	err = r.Request(endpoints.BulkDeleteMessage.Compile(channelID), api.MessageBulkDelete{Messages: messageIDs}, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		// TODO: check here if no err means all messages deleted
		for _, messageID := range messageIDs {
			r.Disgo().Cache().UncacheMessage(channelID, messageID)
		}
	}
	return
}

// CrosspostMessage lets you crosspost a api.Message in a channel with type api.ChannelTypeNews
func (r RestClientImpl) CrosspostMessage(channelID api.Snowflake, messageID api.Snowflake) (msg *api.Message, err error) {
	err = r.Request(endpoints.CrosspostMessage.Compile(channelID, messageID), nil, &msg)
	if err == nil {
		msg = r.Disgo().EntityBuilder().CreateMessage(msg, api.CacheStrategyNoWs)
	}
	return
}

// OpenDMChannel opens a new dm channel a user
func (r RestClientImpl) OpenDMChannel(userID api.Snowflake) (channel *api.DMChannel, err error) {
	body := struct {
		RecipientID api.Snowflake `json:"recipient_id"`
	}{
		RecipientID: userID,
	}
	err = r.Request(endpoints.CreateDMChannel.Compile(), body, &channel)
	if err == nil {
		channel = r.Disgo().EntityBuilder().CreateDMChannel(&channel.MessageChannel.Channel, api.CacheStrategyNoWs)
	}
	return
}

// UpdateSelfNick updates the bots nickname in a guild
func (r RestClientImpl) UpdateSelfNick(guildID api.Snowflake, nick *string) (newNick *string, err error) {
	var updateNick *api.UpdateSelfNick
	err = r.Request(endpoints.UpdateSelfNick.Compile(guildID), api.UpdateSelfNick{Nick: nick}, &updateNick)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().Member(guildID, r.Disgo().SelfUserID()).Nick = updateNick.Nick
		newNick = updateNick.Nick
	}
	return
}

// GetUser fetches the specific user
func (r RestClientImpl) GetUser(userID api.Snowflake) (user *api.User, err error) {
	err = r.Request(endpoints.GetUser.Compile(userID), nil, &user)
	if err == nil {
		user = r.Disgo().EntityBuilder().CreateUser(user, api.CacheStrategyNoWs)
	}
	return
}

// GetMember fetches the specific member
func (r RestClientImpl) GetMember(guildID api.Snowflake, userID api.Snowflake) (member *api.Member, err error) {
	err = r.Request(endpoints.GetMember.Compile(guildID, userID), nil, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// GetMembers fetches all members for a guild
func (r RestClientImpl) GetMembers(guildID api.Snowflake) (members []*api.Member, err error) {
	err = r.Request(endpoints.GetMembers.Compile(guildID), nil, &members)
	if err == nil {
		for _, member := range members {
			member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
		}
	}
	return
}

// AddMember adds a member to the guild with the oauth2 access BotToken. requires api.PermissionCreateInstantInvite
func (r RestClientImpl) AddMember(guildID api.Snowflake, userID api.Snowflake, addGuildMemberData api.AddGuildMemberData) (member *api.Member, err error) {
	err = r.Request(endpoints.AddMember.Compile(guildID, userID), addGuildMemberData, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
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
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMember(guildID, userID)
	}
	return
}

// UpdateMember updates a member
func (r RestClientImpl) UpdateMember(guildID api.Snowflake, userID api.Snowflake, updateGuildMemberData api.UpdateGuildMemberData) (member *api.Member, err error) {
	err = r.Request(endpoints.UpdateMember.Compile(guildID, userID), updateGuildMemberData, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// MoveMember moves/kicks the member to/from a voice channel
func (r RestClientImpl) MoveMember(guildID api.Snowflake, userID api.Snowflake, channelID *api.Snowflake) (member *api.Member, err error) {
	err = r.Request(endpoints.UpdateMember.Compile(guildID, userID), api.MoveGuildMemberData{ChannelID: channelID}, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// AddMemberRole adds a role to a member
func (r RestClientImpl) AddMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
	err = r.Request(endpoints.AddMemberRole.Compile(guildID, userID, roleID), nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			member.Roles = append(member.Roles, roleID)
		}
	}
	return
}

// RemoveMemberRole removes a role from a member
func (r RestClientImpl) RemoveMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
	err = r.Request(endpoints.RemoveMemberRole.Compile(guildID, userID, roleID), nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
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
	if err == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, api.CacheStrategyNoWs)
		}
	}
	return
}

// CreateRole creates a new role for a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) CreateRole(guildID api.Snowflake, role api.UpdateRole) (newRole *api.Role, err error) {
	err = r.Request(endpoints.CreateRole.Compile(guildID), role, &newRole)
	if err == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, api.CacheStrategyNoWs)
	}
	return
}

// UpdateRole updates a role from a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) UpdateRole(guildID api.Snowflake, roleID api.Snowflake, role api.UpdateRole) (newRole *api.Role, err error) {
	err = r.Request(endpoints.UpdateRole.Compile(guildID, roleID), role, &newRole)
	if err == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, api.CacheStrategyNoWs)
	}
	return
}

// UpdateRolePositions updates the position of a role from a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) UpdateRolePositions(guildID api.Snowflake, roleUpdates ...api.UpdateRolePosition) (roles []*api.Role, err error) {
	err = r.Request(endpoints.GetRoles.Compile(guildID), roleUpdates, &roles)
	if err == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, api.CacheStrategyNoWs)
		}
	}
	return
}

// DeleteRole deletes a role from a guild. Requires api.PermissionManageRoles
func (r RestClientImpl) DeleteRole(guildID api.Snowflake, roleID api.Snowflake) (err error) {
	err = r.Request(endpoints.UpdateRole.Compile(guildID, roleID), nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.disgo.Cache().UncacheRole(guildID, roleID)
	}
	return
}

// AddReaction lets you add a reaction to a message_events
func (r RestClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	return r.Request(endpoints.AddReaction.Compile(channelID, messageID, normalizeEmoji(emoji)), nil, nil)
}

// RemoveOwnReaction lets you remove your own reaction from a message_events
func (r RestClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	return r.Request(endpoints.RemoveOwnReaction.Compile(channelID, messageID, normalizeEmoji(emoji)), nil, nil)
}

// RemoveUserReaction lets you remove a specific reaction from a user from a message_events
func (r RestClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) error {
	return r.Request(endpoints.RemoveUserReaction.Compile(channelID, messageID, normalizeEmoji(emoji), userID), nil, nil)
}

// GetGlobalCommands gets you all global commands
func (r RestClientImpl) GetGlobalCommands(applicationID api.Snowflake) (commands []*api.Command, err error) {
	err = r.Request(endpoints.GetGlobalCommands.Compile(applicationID), nil, &commands)
	if err == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGlobalCommand gets you a specific global global command
func (r RestClientImpl) GetGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, err error) {
	err = r.Request(endpoints.GetGlobalCommand.Compile(applicationID, commandID), nil, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// CreateGlobalCommand lets you create a new global command
func (r RestClientImpl) CreateGlobalCommand(applicationID api.Snowflake, command api.Command) (cmd *api.Command, err error) {
	err = r.Request(endpoints.CreateGlobalCommand.Compile(applicationID), command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// SetGlobalCommands lets you override all global commands
func (r RestClientImpl) SetGlobalCommands(applicationID api.Snowflake, commands ...api.Command) (cmds []*api.Command, err error) {
	if len(commands) > 100 {
		err = api.ErrTooMuchApplicationCommands
		return
	}
	err = r.Request(endpoints.SetGlobalCommands.Compile(applicationID), commands, &cmds)
	if err == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// EditGlobalCommand lets you edit a specific global command
func (r RestClientImpl) EditGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.UpdateCommand) (cmd *api.Command, err error) {
	err = r.Request(endpoints.EditGlobalCommand.Compile(applicationID, commandID), command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// DeleteGlobalCommand lets you delete a specific global command
func (r RestClientImpl) DeleteGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (err error) {
	err = r.Request(endpoints.DeleteGlobalCommand.Compile(applicationID, commandID), nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheCommand(commandID)
	}
	return
}

// GetGuildCommands gets you all guild_events commands
func (r RestClientImpl) GetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake) (commands []*api.Command, err error) {
	err = r.Request(endpoints.GetGuildCommands.Compile(applicationID, guildID), nil, &commands)
	if err == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// CreateGuildCommand lets you create a new guild_events command
func (r RestClientImpl) CreateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.Command) (cmd *api.Command, err error) {
	err = r.Request(endpoints.CreateGuildCommand.Compile(applicationID, guildID), command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// SetGuildCommands lets you override all guild_events commands
func (r RestClientImpl) SetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.Command) (cmds []*api.Command, err error) {
	if len(commands) > 100 {
		err = api.ErrTooMuchApplicationCommands
		return
	}
	err = r.Request(endpoints.SetGuildCommands.Compile(applicationID, guildID), commands, &cmds)
	if err == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGuildCommand gets you a specific guild_events command
func (r RestClientImpl) GetGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, err error) {
	err = r.Request(endpoints.GetGuildCommand.Compile(applicationID, guildID, commandID), nil, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// EditGuildCommand lets you edit a specific guild_events command
func (r RestClientImpl) EditGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.UpdateCommand) (cmd *api.Command, err error) {
	err = r.Request(endpoints.EditGuildCommand.Compile(applicationID, guildID, commandID), command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// DeleteGuildCommand lets you delete a specific guild_events command
func (r RestClientImpl) DeleteGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (err error) {
	err = r.Request(endpoints.DeleteGuildCommand.Compile(applicationID, guildID, commandID), nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheCommand(commandID)
	}
	return
}

// GetGuildCommandsPermissions returns the api.CommandPermission for a all api.Command(s) in a guild
func (r RestClientImpl) GetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake) (cmdsPerms []*api.GuildCommandPermissions, err error) {
	err = r.Request(endpoints.GetGuildCommandPermissions.Compile(applicationID, guildID), nil, &cmdsPerms)
	if err == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGuildCommandPermissions returns the api.CommandPermission for a specific api.Command in a guild
func (r RestClientImpl) GetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmdPerms *api.GuildCommandPermissions, err error) {
	err = r.Request(endpoints.GetGuildCommandPermission.Compile(applicationID, guildID, commandID), nil, &cmdPerms)
	if err == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
	}
	return
}

// SetGuildCommandsPermissions sets the api.GuildCommandPermissions for a all api.Command(s)
func (r RestClientImpl) SetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandsPermissions ...api.SetGuildCommandPermissions) (cmdsPerms []*api.GuildCommandPermissions, err error) {
	err = r.Request(endpoints.SetGuildCommandsPermissions.Compile(applicationID, guildID), api.SetGuildCommandsPermissions(commandsPermissions), &cmdsPerms)
	if err == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
		}
	}
	return
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific api.Command
func (r RestClientImpl) SetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, commandPermissions api.SetGuildCommandPermissions) (cmdPerms *api.GuildCommandPermissions, err error) {
	err = r.Request(endpoints.SetGuildCommandPermissions.Compile(applicationID, guildID, commandID), commandPermissions, &cmdPerms)
	if err == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
	}
	return
}

// SendInteractionResponse used to send the initial response on an interaction
func (r RestClientImpl) SendInteractionResponse(interactionID api.Snowflake, interactionToken endpoints.Token, interactionResponse api.InteractionResponse) error {
	return r.Request(endpoints.CreateInteractionResponse.Compile(interactionID, interactionToken), interactionResponse, nil)
}

// EditInteractionResponse used to edit the initial response on an interaction
func (r RestClientImpl) EditInteractionResponse(applicationID api.Snowflake, interactionToken endpoints.Token, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.EditInteractionResponse.Compile(applicationID, interactionToken), followupMessage, &message)
}

// DeleteInteractionResponse used to delete the initial response on an interaction
func (r RestClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken endpoints.Token) error {
	return r.Request(endpoints.DeleteInteractionResponse.Compile(applicationID, interactionToken), nil, nil)
}

// SendFollowupMessage used to send a followup message_events to an interaction
func (r RestClientImpl) SendFollowupMessage(applicationID api.Snowflake, interactionToken endpoints.Token, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.CreateFollowupMessage.Compile(applicationID, interactionToken), followupMessage, &message)
}

// EditFollowupMessage used to edit a api.FollowupMessage from an api.Interaction
func (r RestClientImpl) EditFollowupMessage(applicationID api.Snowflake, interactionToken endpoints.Token, messageID api.Snowflake, followupMessage api.FollowupMessage) (message *api.Message, err error) {
	return message, r.Request(endpoints.EditFollowupMessage.Compile(applicationID, interactionToken, messageID), followupMessage, &message)
}

// DeleteFollowupMessage used to delete a api.FollowupMessage from an api.Interaction
func (r RestClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken endpoints.Token, messageID api.Snowflake) error {
	return r.Request(endpoints.DeleteFollowupMessage.Compile(applicationID, interactionToken, messageID), nil, nil)
}

func normalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}
