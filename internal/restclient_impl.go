package internal

import (
	"net/http"
	"strings"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/restclient"
)

func newRestClientImpl(disgo api.Disgo, httpClient *http.Client) api.RestClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &RestClientImpl{
		RestClient: restclient.NewRestClient(httpClient, disgo.Logger(), "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)", http.Header{"Authorization": []string{"Bot " + disgo.Token()}}),
		disgo:      disgo,
	}
}

// RestClientImpl is the rest client implementation used for HTTP requests to discord
type RestClientImpl struct {
	restclient.RestClient
	disgo api.Disgo
}

func (r *RestClientImpl) GetGateway() (*api.GatewayRs, error) {
	return nil, nil
}

func (r *RestClientImpl) GetGatewayBot() (*api.GatewayBotRs, error) {
	return nil, nil
}

func (r *RestClientImpl) GetBotApplication() (*api.Application, error) {
	return nil, nil
}

func (r *RestClientImpl) GetVoiceRegions() ([]*api.VoiceRegion, error) {
	return nil, nil
}

func (r *RestClientImpl) GetSelfUser() (*api.User, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateSelfUser() (*api.User, error) {
	return nil, nil
}

func (r *RestClientImpl) GetGuilds() ([]*api.Guild, error) {
	return nil, nil
}

func (r *RestClientImpl) LeaveGuild(guildID api.Snowflake) error {
	return nil
}

func (r *RestClientImpl) GetDMChannels() ([]api.DMChannel, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateDMChannel(userID api.Snowflake) (api.DMChannel, error) {
	return nil, nil
}

func (r *RestClientImpl) GetGuild(guildID api.Snowflake) (*api.Guild, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateGuild(guildCreate api.GuildCreate) (*api.Guild, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateGuild(guildUpdate api.GuildUpdate) (*api.Guild, error) {
	return nil, nil
}

func (r *RestClientImpl) DeleteGuild(guildID api.Snowflake) error {
	return nil
}

func (r *RestClientImpl) GetGuildVanityURL(guildID api.Snowflake) (*string, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateGuildChannel(guildID api.Snowflake, channelCreate api.ChannelCreate) (api.GuildChannel, error) {
	return nil, nil
}

func (r *RestClientImpl) GetGuildChannels() ([]api.GuildChannel, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateGuildChannelPositions() error {
	return nil
}

func (r *RestClientImpl) GetBans(guildID api.Snowflake) ([]*api.Ban, error) {
	return nil, nil
}

func (r *RestClientImpl) GetBan(guildID api.Snowflake, userID api.Snowflake) (*api.Ban, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateBan(guildID api.Snowflake, userID api.Snowflake, delDays int, reason string) error {
	return nil
}

func (r *RestClientImpl) DeleteBan(guildID api.Snowflake, userID api.Snowflake) error {
	return nil
}

func (r *RestClientImpl) RemoveMember(guildID api.Snowflake, userID api.Snowflake, reason *string) error {
	return nil
}

func (r *RestClientImpl) GetPruneMembersCount(guildID api.Snowflake, days int, includeRoles []api.Snowflake) (*int, error) {
	return nil, nil
}

func (r *RestClientImpl) PruneMembers(guildID api.Snowflake, days int, computePruneCount bool, includeRoles []api.Snowflake, reason string) (*int, error) {
	return nil, nil
}

func (r *RestClientImpl) GetGuildWebhooks(guildID api.Snowflake) {

}

func (r *RestClientImpl) GetAuditLogs(guildID api.Snowflake) {

}

func (r *RestClientImpl) GetGuildVoiceRegions(guildID api.Snowflake) ([]*api.VoiceRegion, error) {
	return nil, nil
}

func (r *RestClientImpl) GetGuildIntegrations(guildID api.Snowflake) ([]*api.Integration, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateGuildIntegration(guildID api.Snowflake) (*api.Integration, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateGuildIntegration(guildID api.Snowflake) (*api.Integration, error) {
	return nil, nil
}

func (r *RestClientImpl) DeleteGuildIntegration(guildID api.Snowflake) error {
	return nil
}

func (r *RestClientImpl) SyncIntegration(guildID api.Snowflake) {

}

func (r *RestClientImpl) CreateMessage(channelID api.Snowflake, message *api.MessageCreate) (*api.Message, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateMessage(channelID api.Snowflake, messageID api.Snowflake, message *api.MessageUpdate) (*api.Message, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake, command *api.CommandUpdate) (*api.Command, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command *api.CommandUpdate) (*api.Command, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse *api.InteractionResponse) error {
	return nil
}

func (r *RestClientImpl) UpdateInteractionResponse(applicationID api.Snowflake, interactionToken string, followupMessage *api.FollowupMessage) (*api.Message, error) {
	return nil, nil
}

func (r *RestClientImpl) CreateFollowupMessage(applicationID api.Snowflake, interactionToken string, followupMessage *api.FollowupMessage) (*api.Message, error) {
	return nil, nil
}

func (r *RestClientImpl) UpdateFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, followupMessage *api.FollowupMessage) (*api.Message, error) {
	return nil, nil
}

// Disgo returns the api.Disgo instance
func (r *RestClientImpl) Disgo() api.Disgo {
	return r.disgo
}

// Close cleans up the http managers connections
func (r *RestClientImpl) Close() {
	r.HttpClient().CloseIdleConnections()
}

func (r *RestClientImpl) DoWithHeaders(route *restclient.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader http.Header) (err restclient.RestError) {
	err = r.RestClient.DoWithHeaders(route, rqBody, rsBody, customHeader)
	// TODO reimplement events.HTTPRequestEvent 
	/*r.Disgo().EventManager().Dispatch(events.HTTPRequestEvent{
		GenericEvent: events.NewEvent(r.Disgo(), 0),
		Request:      rq,
		Response:     rs,
	}) */

	// TODO reimplement api.ErrorResponse unmarshalling
	/*
		var errorRs api.ErrorResponse
				if err = json.Unmarshal(rawRsBody, &errorRs); err != nil {
					r.Disgo().Logger().Errorf("error unmarshalling error response. code: %d, error: %s", rs.StatusCode, err)
					return err
				}
				return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message_events: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)
	*/
	return
}

// SendMessage lets you send a api.Message to a api.MessageChannel
func (r *RestClientImpl) SendMessage(channelID api.Snowflake, message *api.MessageCreate) (msg *api.Message, err error) {
	compiledRoute, err := restclient.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return nil, err
	}
	var fullMsg *api.FullMessage
	err = r.Do(compiledRoute, message, &fullMsg)
	if err == nil {
		msg = r.Disgo().EntityBuilder().CreateMessage(fullMsg, api.CacheStrategyNoWs)
	}
	return
}

// EditMessage lets you edit a api.Message
func (r *RestClientImpl) EditMessage(channelID api.Snowflake, messageID api.Snowflake, message *api.MessageUpdate) (msg *api.Message, err error) {
	compiledRoute, err := restclient.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, err
	}
	var fullMsg *api.FullMessage
	err = r.Do(compiledRoute, message, &fullMsg)
	if err == nil {
		msg = r.Disgo().EntityBuilder().CreateMessage(fullMsg, api.CacheStrategyNoWs)
	}
	return
}

// DeleteMessage lets you delete a api.Message
func (r *RestClientImpl) DeleteMessage(channelID api.Snowflake, messageID api.Snowflake) (err error) {
	compiledRoute, err := restclient.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMessage(channelID, messageID)
	}
	return
}

// BulkDeleteMessages lets you bulk delete api.Message(s)
func (r *RestClientImpl) BulkDeleteMessages(channelID api.Snowflake, messageIDs ...api.Snowflake) (err error) {
	compiledRoute, err := restclient.BulkDeleteMessage.Compile(nil, channelID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, api.MessageBulkDelete{Messages: messageIDs}, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		// TODO: check here if no err means all messages deleted
		for _, messageID := range messageIDs {
			r.Disgo().Cache().UncacheMessage(channelID, messageID)
		}
	}
	return
}

// CrosspostMessage lets you crosspost a api.Message in a channel with type api.ChannelTypeNews
func (r *RestClientImpl) CrosspostMessage(channelID api.Snowflake, messageID api.Snowflake) (msg *api.Message, err error) {
	compiledRoute, err := restclient.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, err
	}
	var fullMsg *api.FullMessage
	err = r.Do(compiledRoute, nil, &fullMsg)
	if err == nil {
		msg = r.Disgo().EntityBuilder().CreateMessage(fullMsg, api.CacheStrategyNoWs)
	}
	return
}

// OpenDMChannel opens a new dm channel a user
func (r *RestClientImpl) OpenDMChannel(userID api.Snowflake) (dmChannel api.DMChannel, err error) {
	compiledRoute, err := restclient.CreateDMChannel.Compile(nil)
	if err != nil {
		return nil, err
	}
	body := struct {
		RecipientID api.Snowflake `json:"recipient_id"`
	}{
		RecipientID: userID,
	}
	var channel *api.ChannelImpl
	err = r.Do(compiledRoute, body, &channel)
	if err == nil {
		dmChannel = r.Disgo().EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNoWs)
	}
	return
}

// UpdateSelfNick updates the bots nickname in a guild
func (r *RestClientImpl) UpdateSelfNick(guildID api.Snowflake, nick *string) (newNick *string, err error) {
	compiledRoute, err := restclient.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return nil, err
	}
	var updateNick *api.UpdateSelfNick
	err = r.Do(compiledRoute, &api.UpdateSelfNick{Nick: nick}, &updateNick)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().Member(guildID, r.Disgo().ApplicationID()).Nick = updateNick.Nick
		newNick = updateNick.Nick
	}
	return
}

// GetUser fetches the specific user
func (r *RestClientImpl) GetUser(userID api.Snowflake) (user *api.User, err error) {
	compiledRoute, err := restclient.GetUser.Compile(nil, userID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &user)
	if err == nil {
		user = r.Disgo().EntityBuilder().CreateUser(user, api.CacheStrategyNoWs)
	}
	return
}

// GetMember fetches the specific member
func (r *RestClientImpl) GetMember(guildID api.Snowflake, userID api.Snowflake) (member *api.Member, err error) {
	compiledRoute, err := restclient.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// GetMembers fetches all members for a guild
func (r *RestClientImpl) GetMembers(guildID api.Snowflake) (members []*api.Member, err error) {
	compiledRoute, err := restclient.GetMembers.Compile(nil, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &members)
	if err == nil {
		for _, member := range members {
			member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
		}
	}
	return
}

// AddMember adds a member to the guild with the oauth2 access BotToken. requires api.PermissionCreateInstantInvite
func (r *RestClientImpl) AddMember(guildID api.Snowflake, userID api.Snowflake, addGuildMember *api.AddGuildMember) (member *api.Member, err error) {
	compiledRoute, err := restclient.AddMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, addGuildMember, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// KickMember kicks a member from the guild. requires api.PermissionKickMembers
func (r *RestClientImpl) KickMember(guildID api.Snowflake, userID api.Snowflake, reason *string) (err error) {
	var compiledRoute *restclient.CompiledAPIRoute
	var params map[string]interface{}
	if reason != nil {
		params = map[string]interface{}{"reason": *reason}
	}
	compiledRoute, err = restclient.RemoveMember.Compile(params, guildID, userID)
	if err != nil {
		return
	}
	err = r.Do(compiledRoute, nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMember(guildID, userID)
	}
	return
}

// UpdateMember updates a member
func (r *RestClientImpl) UpdateMember(guildID api.Snowflake, userID api.Snowflake, updateGuildMember *api.UpdateGuildMember) (member *api.Member, err error) {
	compiledRoute, err := restclient.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, updateGuildMember, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// MoveMember moves/kicks the member to/from a voice channel
func (r *RestClientImpl) MoveMember(guildID api.Snowflake, userID api.Snowflake, channelID *api.Snowflake) (member *api.Member, err error) {
	compiledRoute, err := restclient.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, api.UpdateGuildMember{ChannelID: channelID}, &member)
	if err == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// AddMemberRole adds a role to a member
func (r *RestClientImpl) AddMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
	compiledRoute, err := restclient.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			member.Roles = append(member.Roles, roleID)
		}
	}
	return
}

// RemoveMemberRole removes a role from a member
func (r *RestClientImpl) RemoveMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
	compiledRoute, err := restclient.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, nil, nil)
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
func (r *RestClientImpl) GetRoles(guildID api.Snowflake) (roles []*api.Role, err error) {
	compiledRoute, err := restclient.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &roles)
	if err == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, api.CacheStrategyNoWs)
		}
	}
	return
}

// CreateRole creates a new role for a guild. Requires api.PermissionManageRoles
func (r *RestClientImpl) CreateRole(guildID api.Snowflake, roleCreate *api.RoleCreate) (newRole *api.Role, err error) {
	compiledRoute, err := restclient.CreateRole.Compile(nil, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, roleCreate, &newRole)
	if err == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, api.CacheStrategyNoWs)
	}
	return
}

// UpdateRole updates a role from a guild. Requires api.PermissionManageRoles
func (r *RestClientImpl) UpdateRole(guildID api.Snowflake, roleID api.Snowflake, roleUpdate *api.RoleUpdate) (newRole *api.Role, err error) {
	compiledRoute, err := restclient.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, roleUpdate, &newRole)
	if err == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, api.CacheStrategyNoWs)
	}
	return
}

// UpdateRolePositions updates the position of a role from a guild. Requires api.PermissionManageRoles
func (r *RestClientImpl) UpdateRolePositions(guildID api.Snowflake, rolePositionUpdates ...*api.RolePositionUpdate) (roles []*api.Role, err error) {
	compiledRoute, err := restclient.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, rolePositionUpdates, &roles)
	if err == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, api.CacheStrategyNoWs)
		}
	}
	return
}

// DeleteRole deletes a role from a guild. Requires api.PermissionManageRoles
func (r *RestClientImpl) DeleteRole(guildID api.Snowflake, roleID api.Snowflake) (err error) {
	compiledRoute, err := restclient.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.disgo.Cache().UncacheRole(guildID, roleID)
	}
	return
}

// AddReaction lets you add a reaction to a message_events
func (r *RestClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	compiledRoute, err := restclient.AddReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji))
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// RemoveOwnReaction lets you remove your own reaction from a message_events
func (r *RestClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	compiledRoute, err := restclient.RemoveOwnReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji))
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// RemoveUserReaction lets you remove a specific reaction from a user from a message_events
func (r *RestClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) error {
	compiledRoute, err := restclient.RemoveUserReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji), userID)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// GetGlobalCommands gets you all global commands
func (r *RestClientImpl) GetGlobalCommands(applicationID api.Snowflake) (commands []*api.Command, err error) {
	compiledRoute, err := restclient.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &commands)
	if err == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGlobalCommand gets you a specific global global command
func (r *RestClientImpl) GetGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, err error) {
	compiledRoute, err := restclient.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// CreateGlobalCommand lets you create a new global command
func (r *RestClientImpl) CreateGlobalCommand(applicationID api.Snowflake, command *api.CommandCreate) (cmd *api.Command, err error) {
	compiledRoute, err := restclient.CreateGlobalCommand.Compile(nil, applicationID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// SetGlobalCommands lets you override all global commands
func (r *RestClientImpl) SetGlobalCommands(applicationID api.Snowflake, commands ...*api.CommandCreate) (cmds []*api.Command, err error) {
	compiledRoute, err := restclient.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, err
	}
	if len(commands) > 100 {
		err = api.ErrTooMuchCommands
		return
	}
	err = r.Do(compiledRoute, commands, &cmds)
	if err == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// EditGlobalCommand lets you edit a specific global command
func (r *RestClientImpl) EditGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake, command *api.CommandUpdate) (cmd *api.Command, err error) {
	compiledRoute, err := restclient.UpdateGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// DeleteGlobalCommand lets you delete a specific global command
func (r *RestClientImpl) DeleteGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (err error) {
	compiledRoute, err := restclient.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheCommand(commandID)
	}
	return
}

// GetGuildCommands gets you all guild_events commands
func (r *RestClientImpl) GetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake) (commands []*api.Command, err error) {
	compiledRoute, err := restclient.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &commands)
	if err == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// CreateGuildCommand lets you create a new guild_events command
func (r *RestClientImpl) CreateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command *api.CommandCreate) (cmd *api.Command, err error) {
	compiledRoute, err := restclient.CreateGuildCommand.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// SetGuildCommands lets you override all guild_events commands
func (r *RestClientImpl) SetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...*api.CommandCreate) (cmds []*api.Command, err error) {
	compiledRoute, err := restclient.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, err
	}
	if len(commands) > 100 {
		err = api.ErrTooMuchCommands
		return
	}
	err = r.Do(compiledRoute, commands, &cmds)
	if err == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGuildCommand gets you a specific guild_events command
func (r *RestClientImpl) GetGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, err error) {
	compiledRoute, err := restclient.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// EditGuildCommand lets you edit a specific guild_events command
func (r *RestClientImpl) EditGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command *api.CommandUpdate) (cmd *api.Command, err error) {
	compiledRoute, err := restclient.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, command, &cmd)
	if err == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// DeleteGuildCommand lets you delete a specific guild_events command
func (r *RestClientImpl) DeleteGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (err error) {
	compiledRoute, err := restclient.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return err
	}
	err = r.Do(compiledRoute, nil, nil)
	if err == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheCommand(commandID)
	}
	return
}

// GetGuildCommandsPermissions returns the api.CommandPermission for a all api.Command(s) in a guild
func (r *RestClientImpl) GetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake) (cmdsPerms []*api.GuildCommandPermissions, err error) {
	compiledRoute, err := restclient.GetGuildCommandPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &cmdsPerms)
	if err == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGuildCommandPermissions returns the api.CommandPermission for a specific api.Command in a guild
func (r *RestClientImpl) GetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmdPerms *api.GuildCommandPermissions, err error) {
	compiledRoute, err := restclient.GetGuildCommandPermission.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, nil, &cmdPerms)
	if err == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
	}
	return
}

// SetGuildCommandsPermissions sets the api.GuildCommandPermissions for a all api.Command(s)
func (r *RestClientImpl) SetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandsPermissions ...*api.SetGuildCommandPermissions) (cmdsPerms []*api.GuildCommandPermissions, err error) {
	compiledRoute, err := restclient.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, api.SetGuildCommandsPermissions(commandsPermissions), &cmdsPerms)
	if err == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
		}
	}
	return
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific api.Command
func (r *RestClientImpl) SetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, commandPermissions *api.SetGuildCommandPermissions) (cmdPerms *api.GuildCommandPermissions, err error) {
	compiledRoute, err := restclient.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, err
	}
	err = r.Do(compiledRoute, commandPermissions, &cmdPerms)
	if err == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
	}
	return
}

// SendInteractionResponse used to send the initial response on an interaction
func (r *RestClientImpl) SendInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse *api.InteractionResponse) error {
	compiledRoute, err := restclient.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, interactionResponse, nil)
}

// EditInteractionResponse used to edit the initial response on an interaction
func (r *RestClientImpl) EditInteractionResponse(applicationID api.Snowflake, interactionToken string, followupMessage *api.FollowupMessage) (message *api.Message, err error) {
	compiledRoute, err := restclient.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, err
	}
	return message, r.Do(compiledRoute, followupMessage, &message)
}

// DeleteInteractionResponse used to delete the initial response on an interaction
func (r *RestClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) error {
	compiledRoute, err := restclient.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// SendFollowupMessage used to send a followup message_events to an interaction
func (r *RestClientImpl) SendFollowupMessage(applicationID api.Snowflake, interactionToken string, followupMessage *api.FollowupMessage) (message *api.Message, err error) {
	compiledRoute, err := restclient.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, err
	}
	return message, r.Do(compiledRoute, followupMessage, &message)
}

// EditFollowupMessage used to edit a api.FollowupMessage from an api.Interaction
func (r *RestClientImpl) EditFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, followupMessage *api.FollowupMessage) (message *api.Message, err error) {
	compiledRoute, err := restclient.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return nil, err
	}
	return message, r.Do(compiledRoute, followupMessage, &message)
}

// DeleteFollowupMessage used to delete a api.FollowupMessage from an api.Interaction
func (r *RestClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) error {
	compiledRoute, err := restclient.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

func normalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}
