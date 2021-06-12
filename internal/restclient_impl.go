package internal

import (
	"net/http"
	"strings"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/restclient"
)

var _ api.RestClient = (*restClientImpl)(nil)

func newRestClientImpl(disgo api.Disgo, httpClient *http.Client) api.RestClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &restClientImpl{
		RestClient: restclient.NewRestClient(httpClient, disgo.Logger(), api.UserAgent, http.Header{"Authorization": []string{"Bot " + disgo.Token()}}),
		disgo:      disgo,
	}
}

// restClientImpl is the rest client implementation used for HTTP requests to discord
type restClientImpl struct {
	restclient.RestClient
	disgo api.Disgo
}

func (r *restClientImpl) GetSelfUser() (*api.User, error) {
	panic("implement me")
}

func (r *restClientImpl) UpdateSelfUser() (*api.User, error) {
	panic("implement me")
}

func (r *restClientImpl) GetGuilds() ([]*api.Guild, error) {
	panic("implement me")
}

func (r *restClientImpl) LeaveGuild(guildID api.Snowflake) error {
	panic("implement me")
}

func (r *restClientImpl) GetDMChannels() ([]api.DMChannel, error) {
	panic("implement me")
}

func (r *restClientImpl) GetGuild(guildID api.Snowflake) (*api.Guild, error) {
	panic("implement me")
}

func (r *restClientImpl) CreateGuild(guildCreate api.GuildCreate) (*api.Guild, error) {
	panic("implement me")
}

func (r *restClientImpl) UpdateGuild(guildUpdate api.GuildUpdate) (*api.Guild, error) {
	panic("implement me")
}

func (r *restClientImpl) DeleteGuild(guildID api.Snowflake) error {
	panic("implement me")
}

func (r *restClientImpl) GetGuildVanityURL(guildID api.Snowflake) (*string, error) {
	panic("implement me")
}

func (r *restClientImpl) CreateGuildChannel(guildID api.Snowflake, channelCreate api.ChannelCreate) (api.GuildChannel, error) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildChannels() ([]api.GuildChannel, error) {
	panic("implement me")
}

func (r *restClientImpl) UpdateGuildChannelPositions() error {
	panic("implement me")
}

func (r *restClientImpl) GetBans(guildID api.Snowflake) ([]*api.Ban, error) {
	panic("implement me")
}

func (r *restClientImpl) GetBan(guildID api.Snowflake, userID api.Snowflake) (*api.Ban, error) {
	panic("implement me")
}

func (r *restClientImpl) CreateBan(guildID api.Snowflake, userID api.Snowflake, delDays int, reason string) error {
	panic("implement me")
}

func (r *restClientImpl) DeleteBan(guildID api.Snowflake, userID api.Snowflake) error {
	panic("implement me")
}

func (r *restClientImpl) GetPruneMembersCount(guildID api.Snowflake, days int, includeRoles []api.Snowflake) (*int, error) {
	panic("implement me")
}

func (r *restClientImpl) PruneMembers(guildID api.Snowflake, days int, computePruneCount bool, includeRoles []api.Snowflake, reason string) (*int, error) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildWebhooks(guildID api.Snowflake) {
	panic("implement me")
}

func (r *restClientImpl) GetAuditLogs(guildID api.Snowflake) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildVoiceRegions(guildID api.Snowflake) ([]*api.VoiceRegion, error) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildIntegrations(guildID api.Snowflake) ([]*api.Integration, error) {
	panic("implement me")
}

func (r *restClientImpl) CreateGuildIntegration(guildID api.Snowflake) (*api.Integration, error) {
	panic("implement me")
}

func (r *restClientImpl) UpdateGuildIntegration(guildID api.Snowflake) (*api.Integration, error) {
	panic("implement me")
}

func (r *restClientImpl) DeleteGuildIntegration(guildID api.Snowflake) error {
	panic("implement me")
}

func (r *restClientImpl) SyncIntegration(guildID api.Snowflake) {
	panic("implement me")
}

// Disgo returns the api.Disgo instance
func (r *restClientImpl) Disgo() api.Disgo {
	return r.disgo
}

// Close cleans up the http managers connections
func (r *restClientImpl) Close() {
	r.HTTPClient().CloseIdleConnections()
}

// DoWithHeaders executes a rest request with custom headers
func (r *restClientImpl) DoWithHeaders(route *restclient.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader http.Header) (err restclient.RestError) {
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

// CreateMessage lets you send a api.Message to a api.MessageChannel
func (r *restClientImpl) CreateMessage(channelID api.Snowflake, messageCreate api.MessageCreate) (message *api.Message, err error) {
	compiledRoute, err := restclient.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return nil, err
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return nil, err
	}

	var fullMessage *api.FullMessage
	err = r.Do(compiledRoute, body, &fullMessage)
	if err == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(fullMessage, api.CacheStrategyNoWs)
	}
	return
}

// UpdateMessage lets you edit a api.Message
func (r *restClientImpl) UpdateMessage(channelID api.Snowflake, messageID api.Snowflake, messageUpdate api.MessageUpdate) (message *api.Message, err error) {
	compiledRoute, err := restclient.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, err
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, err
	}

	var fullMessage *api.FullMessage
	err = r.Do(compiledRoute, body, &fullMessage)
	if err == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(fullMessage, api.CacheStrategyNoWs)
	}
	return
}

// DeleteMessage lets you delete a api.Message
func (r *restClientImpl) DeleteMessage(channelID api.Snowflake, messageID api.Snowflake) (err error) {
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
func (r *restClientImpl) BulkDeleteMessages(channelID api.Snowflake, messageIDs ...api.Snowflake) (err error) {
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
func (r *restClientImpl) CrosspostMessage(channelID api.Snowflake, messageID api.Snowflake) (msg *api.Message, err error) {
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

// CreateDMChannel opens a new dm channel a user
func (r *restClientImpl) CreateDMChannel(userID api.Snowflake) (dmChannel api.DMChannel, err error) {
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
func (r *restClientImpl) UpdateSelfNick(guildID api.Snowflake, nick *string) (newNick *string, err error) {
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
func (r *restClientImpl) GetUser(userID api.Snowflake) (user *api.User, err error) {
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
func (r *restClientImpl) GetMember(guildID api.Snowflake, userID api.Snowflake) (member *api.Member, err error) {
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
func (r *restClientImpl) GetMembers(guildID api.Snowflake) (members []*api.Member, err error) {
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
func (r *restClientImpl) AddMember(guildID api.Snowflake, userID api.Snowflake, addGuildMember api.AddGuildMember) (member *api.Member, err error) {
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

// KickMember kicks a api.Member from the api.Guild. requires api.PermissionKickMembers
func (r *restClientImpl) KickMember(guildID api.Snowflake, userID api.Snowflake, reason *string) (err error) {
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

// UpdateMember updates a api.MemberUpdateGuildMember
func (r *restClientImpl) UpdateMember(guildID api.Snowflake, userID api.Snowflake, updateGuildMember api.UpdateGuildMember) (member *api.Member, err error) {
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

// MoveMember moves/kicks the api.Member to/from a api.VoiceChannel
func (r *restClientImpl) MoveMember(guildID api.Snowflake, userID api.Snowflake, channelID *api.Snowflake) (member *api.Member, err error) {
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

// AddMemberRole adds a api.Role to a api.Member
func (r *restClientImpl) AddMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
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

// RemoveMemberRole removes a api.Role(s) from a api.Member
func (r *restClientImpl) RemoveMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (err error) {
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

// GetRoles fetches all api.Role(s) from a api.Guild
func (r *restClientImpl) GetRoles(guildID api.Snowflake) (roles []*api.Role, err error) {
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
func (r *restClientImpl) CreateRole(guildID api.Snowflake, roleCreate api.RoleCreate) (newRole *api.Role, err error) {
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
func (r *restClientImpl) UpdateRole(guildID api.Snowflake, roleID api.Snowflake, roleUpdate api.RoleUpdate) (newRole *api.Role, err error) {
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
func (r *restClientImpl) UpdateRolePositions(guildID api.Snowflake, rolePositionUpdates ...api.RolePositionUpdate) (roles []*api.Role, err error) {
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
func (r *restClientImpl) DeleteRole(guildID api.Snowflake, roleID api.Snowflake) (err error) {
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

// AddReaction lets you add a reaction to a api.Message
func (r *restClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	compiledRoute, err := restclient.AddReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji))
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// RemoveOwnReaction lets you remove your own reaction from a api.Message
func (r *restClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) error {
	compiledRoute, err := restclient.RemoveOwnReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji))
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// RemoveUserReaction lets you remove a specific reaction from a api.User from a api.Message
func (r *restClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) error {
	compiledRoute, err := restclient.RemoveUserReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji), userID)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// GetGlobalCommands gets you all global api.Command(s)
func (r *restClientImpl) GetGlobalCommands(applicationID api.Snowflake) (commands []*api.Command, err error) {
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

// GetGlobalCommand gets you a specific global global api.Command
func (r *restClientImpl) GetGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, err error) {
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

// CreateGlobalCommand lets you create a new global api.Command
func (r *restClientImpl) CreateGlobalCommand(applicationID api.Snowflake, command api.CommandCreate) (cmd *api.Command, err error) {
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

// SetGlobalCommands lets you override all global api.Command
func (r *restClientImpl) SetGlobalCommands(applicationID api.Snowflake, commands ...api.CommandCreate) (cmds []*api.Command, err error) {
	compiledRoute, err := restclient.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, err
	}
	if len(commands) > 100 {
		err = api.ErrMaxCommands
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

// UpdateGlobalCommand lets you edit a specific global api.Command
func (r *restClientImpl) UpdateGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.CommandUpdate) (cmd *api.Command, err error) {
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

// DeleteGlobalCommand lets you delete a specific global api.Command
func (r *restClientImpl) DeleteGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (err error) {
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

// GetGuildCommands gets you all api.Command(s) from a api.Guild
func (r *restClientImpl) GetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake) (commands []*api.Command, err error) {
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

// CreateGuildCommand lets you create a new api.Command in a api.Guild
func (r *restClientImpl) CreateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.CommandCreate) (cmd *api.Command, err error) {
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

// SetGuildCommands lets you override all api.Command(s) in a api.Guild
func (r *restClientImpl) SetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.CommandCreate) (cmds []*api.Command, err error) {
	compiledRoute, err := restclient.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, err
	}
	if len(commands) > 100 {
		err = api.ErrMaxCommands
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

// GetGuildCommand gets you a specific api.Command in a api.Guild
func (r *restClientImpl) GetGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, err error) {
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

// UpdateGuildCommand lets you edit a specific api.Command in a api.Guild
func (r *restClientImpl) UpdateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.CommandUpdate) (cmd *api.Command, err error) {
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

// DeleteGuildCommand lets you delete a specific api.Command in a api.Guild
func (r *restClientImpl) DeleteGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (err error) {
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

// GetGuildCommandsPermissions returns the api.CommandPermission for a all api.Command(s) in a api.Guild
func (r *restClientImpl) GetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake) (cmdsPerms []*api.GuildCommandPermissions, err error) {
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

// GetGuildCommandPermissions returns the api.CommandPermission for a specific api.Command in a api.Guild
func (r *restClientImpl) GetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmdPerms *api.GuildCommandPermissions, err error) {
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
func (r *restClientImpl) SetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandsPermissions ...api.SetGuildCommandPermissions) (cmdsPerms []*api.GuildCommandPermissions, err error) {
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
func (r *restClientImpl) SetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, commandPermissions api.SetGuildCommandPermissions) (cmdPerms *api.GuildCommandPermissions, err error) {
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

// CreateInteractionResponse used to send the initial response on an api.Interaction
func (r *restClientImpl) CreateInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) error {
	compiledRoute, err := restclient.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return err
	}

	body, err := interactionResponse.ToBody()
	if err != nil {
		return err
	}

	return r.Do(compiledRoute, body, nil)
}

// UpdateInteractionResponse used to edit the initial response on an api.Interaction
func (r *restClientImpl) UpdateInteractionResponse(applicationID api.Snowflake, interactionToken string, messageUpdate api.MessageUpdate) (message *api.Message, err error) {
	compiledRoute, err := restclient.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, err
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, err
	}

	var fullMessage *api.FullMessage
	err = r.Do(compiledRoute, body, &fullMessage)
	if err == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(fullMessage, api.CacheStrategyNoWs)
	}
	return
}

// DeleteInteractionResponse used to delete the initial response on an api.Interaction
func (r *restClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) error {
	compiledRoute, err := restclient.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

// CreateFollowupMessage used to send a followup api.Message to an api.Interaction
func (r *restClientImpl) CreateFollowupMessage(applicationID api.Snowflake, interactionToken string, messageCreate api.MessageCreate) (message *api.Message, err error) {
	compiledRoute, err := restclient.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, err
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return nil, err
	}

	var fullMessage *api.FullMessage
	err = r.Do(compiledRoute, body, &fullMessage)
	if err == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(fullMessage, api.CacheStrategyNoWs)
	}

	return
}

// UpdateFollowupMessage used to edit a followup api.Message from an api.Interaction
func (r *restClientImpl) UpdateFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, messageUpdate api.MessageUpdate) (message *api.Message, err error) {
	compiledRoute, err := restclient.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return nil, err
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, err
	}

	var fullMessage *api.FullMessage
	err = r.Do(compiledRoute, body, &fullMessage)
	if err == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(fullMessage, api.CacheStrategyNoWs)
	}

	return
}

// DeleteFollowupMessage used to delete a followup api.Message from an api.Interaction
func (r *restClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) error {
	compiledRoute, err := restclient.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return err
	}
	return r.Do(compiledRoute, nil, nil)
}

func normalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}
