package internal

import (
	"net/http"
	"time"

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

func (r *restClientImpl) CreateThreadWithMessage(channelID api.Snowflake, messageID api.Snowflake, threadCreate api.ThreadCreate) (api.Thread, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) CreateThreadWithoutMessage(channelID api.Snowflake, threadCreate api.ThreadCreate) (api.Thread, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) JoinThread(threadID api.Snowflake) restclient.RestError {
	panic("implement me")
}

func (r *restClientImpl) AddThreadMember(threadID api.Snowflake, userID api.Snowflake) restclient.RestError {
	panic("implement me")
}

func (r *restClientImpl) LeaveThread(threadID api.Snowflake) restclient.RestError {
	panic("implement me")
}

func (r *restClientImpl) RemoveThreadMember(threadID api.Snowflake, userID api.Snowflake) restclient.RestError {
	panic("implement me")
}

func (r *restClientImpl) GetThreadMembers(threadID api.Snowflake) ([]*api.ThreadMember, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetActiveThreads(channelID api.Snowflake) ([]api.Thread, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetPublicArchivedThreads(channelID api.Snowflake, before time.Time, limit int) ([]api.Thread, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetPrivateArchivedThreads(channelID api.Snowflake, before time.Time, limit int) ([]api.Thread, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetJoinedPrivateArchivedThreads(channelID api.Snowflake, before time.Time, limit int) ([]api.Thread, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetPruneMembersCount(guildID api.Snowflake, days int, includeRoles []api.Snowflake) (*int, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) PruneMembers(guildID api.Snowflake, days int, computePruneCount bool, includeRoles []api.Snowflake, reason string) (*int, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildWebhooks(guildID api.Snowflake) {
	panic("implement me")
}

func (r *restClientImpl) GetAuditLogs(guildID api.Snowflake) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildVoiceRegions(guildID api.Snowflake) ([]*api.VoiceRegion, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) GetGuildIntegrations(guildID api.Snowflake) ([]*api.Integration, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) CreateGuildIntegration(guildID api.Snowflake) (*api.Integration, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) UpdateGuildIntegration(guildID api.Snowflake) (*api.Integration, restclient.RestError) {
	panic("implement me")
}

func (r *restClientImpl) DeleteGuildIntegration(guildID api.Snowflake) restclient.RestError {
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
func (r *restClientImpl) DoWithHeaders(route *restclient.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader http.Header) (rErr restclient.RestError) {
	err := r.RestClient.DoWithHeaders(route, rqBody, rsBody, customHeader)
	rErr = restclient.NewError(nil, err)
	// TODO reimplement events.HTTPRequestEvent
	/*r.Disgo().EventManager().Dispatch(&events.HTTPRequestEvent{
		GenericEvent: events.NewGenericEvent(r.Disgo(), 0),
		Request:      rq,
		Response:     rs,
	}) */

	// TODO reimplement api.ErrorResponse unmarshalling
	/*
		var errorRs api.ErrorResponse
				if err = json.Unmarshal(rawRsBody, &errorRs); err != nil {
					r.Disgo().Logger().Errorf("restclient.RestError unmarshalling restclient.RestError response. code: %d, restclient.RestError: %s", rs.StatusCode, err)
					return err
				}
				return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message_events: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)
	*/
	return
}

// GetUser fetches the specific user
func (r *restClientImpl) GetUser(userID api.Snowflake) (user *api.User, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetUser.Compile(nil, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &user)
	if rErr == nil {
		user = r.Disgo().EntityBuilder().CreateUser(user, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetSelfUser() (selfUser *api.SelfUser, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	var user *api.User
	rErr = r.Do(compiledRoute, nil, &user)
	if rErr == nil {
		selfUser = &api.SelfUser{User: r.Disgo().EntityBuilder().CreateUser(user, api.CacheStrategyNoWs)}
	}
	return
}

func (r *restClientImpl) UpdateSelfUser(updateSelfUser api.UpdateSelfUser) (selfUser *api.SelfUser, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	var user *api.User
	rErr = r.Do(compiledRoute, updateSelfUser, &user)
	if rErr == nil {
		selfUser = &api.SelfUser{User: r.Disgo().EntityBuilder().CreateUser(user, api.CacheStrategyNoWs)}
	}
	return
}

func (r *restClientImpl) GetGuilds(before int, after int, limit int) (guilds []*api.PartialGuild, rErr restclient.RestError) {
	queryParams := restclient.QueryValues{}
	if before > 0 {
		queryParams["before"] = before
	}
	if after > 0 {
		queryParams["after"] = after
	}
	if limit > 0 {
		queryParams["limit"] = limit
	}
	compiledRoute, err := restclient.GetGuilds.Compile(queryParams)
	if err != nil {
		return nil, restclient.NewError(nil, restclient.NewError(nil, err))
	}

	rErr = r.Do(compiledRoute, nil, &guilds)
	return
}

func (r *restClientImpl) LeaveGuild(guildID api.Snowflake) restclient.RestError {
	compiledRoute, err := restclient.LeaveGuild.Compile(nil, guildID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

func (r *restClientImpl) GetDMChannels() (dmChannels []api.DMChannel, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetDMChannels.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var channels []*api.ChannelImpl
	rErr = r.Do(compiledRoute, nil, &channels)
	if rErr == nil {
		dmChannels = make([]api.DMChannel, len(channels))
		for i, channel := range channels {
			dmChannels[i] = r.Disgo().EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNoWs)
		}
	}
	return
}

func (r *restClientImpl) GetMessage(channelID api.Snowflake, messageID api.Snowflake) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}
	return
}

// CreateMessage lets you send a api.Message to a api.MessageChannel
func (r *restClientImpl) CreateMessage(channelID api.Snowflake, messageCreate api.MessageCreate) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, body, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}
	return
}

// UpdateMessage lets you edit a api.Message
func (r *restClientImpl) UpdateMessage(channelID api.Snowflake, messageID api.Snowflake, messageUpdate api.MessageUpdate) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, body, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}
	return
}

// DeleteMessage lets you delete a api.Message
func (r *restClientImpl) DeleteMessage(channelID api.Snowflake, messageID api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMessage(channelID, messageID)
	}
	return
}

// BulkDeleteMessages lets you bulk delete api.Message(s)
func (r *restClientImpl) BulkDeleteMessages(channelID api.Snowflake, messageIDs ...api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.BulkDeleteMessage.Compile(nil, channelID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, api.MessageBulkDelete{Messages: messageIDs}, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		// TODO: check here if no err means all messages deleted
		for _, messageID := range messageIDs {
			r.Disgo().Cache().UncacheMessage(channelID, messageID)
		}
	}
	return
}

// CrosspostMessage lets you crosspost a api.Message in a channel with type api.ChannelTypeNews
func (r *restClientImpl) CrosspostMessage(channelID api.Snowflake, messageID api.Snowflake) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}
	return
}

// CreateDMChannel opens a new api.DMChannel a user
func (r *restClientImpl) CreateDMChannel(userID api.Snowflake) (dmChannel api.DMChannel, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateDMChannel.Compile(nil, struct{ Recipient api.Snowflake }{Recipient: userID})
	if err != nil {
		return nil, restclient.NewError(nil, restclient.NewError(nil, err))
	}

	var channel *api.ChannelImpl
	rErr = r.Do(compiledRoute, dmChannel, &channel)
	if rErr == nil {
		dmChannel = r.Disgo().EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetGuild(guildID api.Snowflake, withCounts bool) (guild *api.Guild, rErr restclient.RestError) {
	var queryParams = restclient.QueryValues{
		"with_counts": withCounts,
	}
	compiledRoute, err := restclient.GetGuild.Compile(queryParams, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, restclient.NewError(nil, err))
	}

	var fullGuild *api.FullGuild
	rErr = r.Do(compiledRoute, nil, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetGuildPreview(guildID api.Snowflake) (guildPreview *api.GuildPreview, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildPreview.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, restclient.NewError(nil, err))
	}

	rErr = r.Do(compiledRoute, nil, &guildPreview)
	return
}

func (r *restClientImpl) CreateGuild(guildCreate api.GuildCreate) (guild *api.Guild, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateGuild.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var fullGuild *api.FullGuild
	rErr = r.Do(compiledRoute, guildCreate, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) UpdateGuild(guildID api.Snowflake, guildUpdate api.GuildUpdate) (guild *api.Guild, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateGuild.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var fullGuild *api.FullGuild
	rErr = r.Do(compiledRoute, guildUpdate, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) DeleteGuild(guildID api.Snowflake) restclient.RestError {
	compiledRoute, err := restclient.DeleteGuild.Compile(nil, guildID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

// GetMember fetches the specific member
func (r *restClientImpl) GetMember(guildID api.Snowflake, userID api.Snowflake) (member *api.Member, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// GetMembers fetches all members for a guild
func (r *restClientImpl) GetMembers(guildID api.Snowflake) (members []*api.Member, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetMembers.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &members)
	if rErr == nil {
		for _, member := range members {
			member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
		}
	}
	return
}

func (r *restClientImpl) SearchMembers(guildID api.Snowflake, query string, limit int) (members []*api.Member, rErr restclient.RestError) {
	queryParams := restclient.QueryValues{}
	if query != "" {
		queryParams["query"] = query
	}
	if limit > 0 {
		queryParams["limit"] = limit
	}
	compiledRoute, err := restclient.GetMembers.Compile(queryParams, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &members)
	if rErr == nil {
		members = make([]*api.Member, len(members))
		for i, member := range members {
			members[i] = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
		}
	}
	return
}

// AddMember adds a member to the guild with the oauth2 access BotToken. requires api.PermissionCreateInstantInvite
func (r *restClientImpl) AddMember(guildID api.Snowflake, userID api.Snowflake, memberAdd api.MemberAdd) (member *api.Member, rErr restclient.RestError) {
	compiledRoute, err := restclient.AddMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, memberAdd, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// RemoveMember kicks a api.Member from the api.Guild. requires api.PermissionKickMembers
func (r *restClientImpl) RemoveMember(guildID api.Snowflake, userID api.Snowflake, reason string) (rErr restclient.RestError) {
	var params restclient.QueryValues
	if reason != "" {
		params = restclient.QueryValues{"reason": reason}
	}
	compiledRoute, err := restclient.RemoveMember.Compile(params, guildID, userID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMember(guildID, userID)
	}
	return
}

// UpdateMember updates a api.Member
func (r *restClientImpl) UpdateMember(guildID api.Snowflake, userID api.Snowflake, memberUpdate api.MemberUpdate) (member *api.Member, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, memberUpdate, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// UpdateSelfNick updates the bots nickname in a guild
func (r *restClientImpl) UpdateSelfNick(guildID api.Snowflake, nick string) (newNick *string, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	var updateNick *api.UpdateSelfNick
	rErr = r.Do(compiledRoute, &api.UpdateSelfNick{Nick: nick}, &updateNick)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		var nick *string
		if updateNick.Nick == "" {
			nick = nil
		}
		r.Disgo().Cache().Member(guildID, r.Disgo().ClientID()).Nick = nick
		newNick = nick
	}
	return
}

// MoveMember moves/kicks the api.Member to/from a api.VoiceChannel
func (r *restClientImpl) MoveMember(guildID api.Snowflake, userID api.Snowflake, channelID *api.Snowflake) (member *api.Member, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, api.MemberMove{ChannelID: channelID}, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, api.CacheStrategyNoWs)
	}
	return
}

// AddMemberRole adds a api.Role to a api.Member
func (r *restClientImpl) AddMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			member.RoleIDs = append(member.RoleIDs, roleID)
		}
	}
	return
}

// RemoveMemberRole removes a api.Role(s) from a api.Member
func (r *restClientImpl) RemoveMemberRole(guildID api.Snowflake, userID api.Snowflake, roleID api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			for i, id := range member.RoleIDs {
				if id == roleID {
					member.RoleIDs = append(member.RoleIDs[:i], member.RoleIDs[i+1:]...)
					break
				}
			}
		}
	}
	return
}

// GetRoles fetches all api.Role(s) from a api.Guild
func (r *restClientImpl) GetRoles(guildID api.Snowflake) (roles []*api.Role, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &roles)
	if rErr == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, api.CacheStrategyNoWs)
		}
	}
	return
}

// CreateRole creates a new role for a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) CreateRole(guildID api.Snowflake, roleCreate api.RoleCreate) (newRole *api.Role, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateRole.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, roleCreate, &newRole)
	if rErr == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, api.CacheStrategyNoWs)
	}
	return
}

// UpdateRole updates a role from a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) UpdateRole(guildID api.Snowflake, roleID api.Snowflake, roleUpdate api.RoleUpdate) (newRole *api.Role, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, roleUpdate, &newRole)
	if rErr == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, api.CacheStrategyNoWs)
	}
	return
}

// UpdateRolePositions updates the position of a role from a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) UpdateRolePositions(guildID api.Snowflake, rolePositionUpdates ...api.RolePositionUpdate) (roles []*api.Role, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, rolePositionUpdates, &roles)
	if rErr == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, api.CacheStrategyNoWs)
		}
	}
	return
}

// DeleteRole deletes a role from a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) DeleteRole(guildID api.Snowflake, roleID api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.disgo.Cache().UncacheRole(guildID, roleID)
	}
	return
}

// AddReaction lets you add a reaction to a api.Message
func (r *restClientImpl) AddReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) restclient.RestError {
	compiledRoute, err := restclient.AddReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji))
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

// RemoveOwnReaction lets you remove your own reaction from a api.Message
func (r *restClientImpl) RemoveOwnReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string) restclient.RestError {
	compiledRoute, err := restclient.RemoveOwnReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji))
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

// RemoveUserReaction lets you remove a specific reaction from a api.User from a api.Message
func (r *restClientImpl) RemoveUserReaction(channelID api.Snowflake, messageID api.Snowflake, emoji string, userID api.Snowflake) restclient.RestError {
	compiledRoute, err := restclient.RemoveUserReaction.Compile(nil, channelID, messageID, normalizeEmoji(emoji), userID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

// GetGlobalCommands gets you all global api.Command(s)
func (r *restClientImpl) GetGlobalCommands(applicationID api.Snowflake) (commands []*api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &commands)
	if rErr == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGlobalCommand gets you a specific global global api.Command
func (r *restClientImpl) GetGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// CreateGlobalCommand lets you create a new global api.Command
func (r *restClientImpl) CreateGlobalCommand(applicationID api.Snowflake, command api.CommandCreate) (cmd *api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateGlobalCommand.Compile(nil, applicationID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// SetGlobalCommands lets you override all global api.Command
func (r *restClientImpl) SetGlobalCommands(applicationID api.Snowflake, commands ...api.CommandCreate) (cmds []*api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	if len(commands) > 100 {
		err = api.ErrMaxCommands
		return
	}
	rErr = r.Do(compiledRoute, commands, &cmds)
	if rErr == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// UpdateGlobalCommand lets you edit a specific global api.Command
func (r *restClientImpl) UpdateGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake, command api.CommandUpdate) (cmd *api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, api.CacheStrategyNoWs)
	}
	return
}

// DeleteGlobalCommand lets you delete a specific global api.Command
func (r *restClientImpl) DeleteGlobalCommand(applicationID api.Snowflake, commandID api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheCommand(commandID)
	}
	return
}

// GetGuildCommands gets you all api.Command(s) from a api.Guild
func (r *restClientImpl) GetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake) (commands []*api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &commands)
	if rErr == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// CreateGuildCommand lets you create a new api.Command in a api.Guild
func (r *restClientImpl) CreateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, command api.CommandCreate) (cmd *api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateGuildCommand.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// SetGuildCommands lets you override all api.Command(s) in a api.Guild
func (r *restClientImpl) SetGuildCommands(applicationID api.Snowflake, guildID api.Snowflake, commands ...api.CommandCreate) (cmds []*api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	if len(commands) > 100 {
		err = api.ErrMaxCommands
		return
	}
	rErr = r.Do(compiledRoute, commands, &cmds)
	if rErr == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGuildCommand gets you a specific api.Command in a api.Guild
func (r *restClientImpl) GetGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmd *api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// UpdateGuildCommand lets you edit a specific api.Command in a api.Guild
func (r *restClientImpl) UpdateGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, command api.CommandUpdate) (cmd *api.Command, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, api.CacheStrategyNoWs)
	}
	return
}

// DeleteGuildCommand lets you delete a specific api.Command in a api.Guild
func (r *restClientImpl) DeleteGuildCommand(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (rErr restclient.RestError) {
	compiledRoute, err := restclient.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && api.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheCommand(commandID)
	}
	return
}

// GetGuildCommandsPermissions returns the api.CommandPermission for a all api.Command(s) in a api.Guild
func (r *restClientImpl) GetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake) (cmdsPerms []*api.GuildCommandPermissions, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildCommandPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmdsPerms)
	if rErr == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
		}
	}
	return
}

// GetGuildCommandPermissions returns the api.CommandPermission for a specific api.Command in a api.Guild
func (r *restClientImpl) GetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake) (cmdPerms *api.GuildCommandPermissions, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildCommandPermission.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmdPerms)
	if rErr == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
	}
	return
}

// SetGuildCommandsPermissions sets the api.GuildCommandPermissions for a all api.Command(s)
func (r *restClientImpl) SetGuildCommandsPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandsPermissions ...api.SetGuildCommandPermissions) (cmdsPerms []*api.GuildCommandPermissions, rErr restclient.RestError) {
	compiledRoute, err := restclient.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, api.SetGuildCommandsPermissions(commandsPermissions), &cmdsPerms)
	if rErr == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
		}
	}
	return
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific api.Command
func (r *restClientImpl) SetGuildCommandPermissions(applicationID api.Snowflake, guildID api.Snowflake, commandID api.Snowflake, commandPermissions api.SetGuildCommandPermissions) (cmdPerms *api.GuildCommandPermissions, rErr restclient.RestError) {
	compiledRoute, err := restclient.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, commandPermissions, &cmdPerms)
	if rErr == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, api.CacheStrategyNoWs)
	}
	return
}

// CreateInteractionResponse used to send the initial response on an api.Interaction
func (r *restClientImpl) CreateInteractionResponse(interactionID api.Snowflake, interactionToken string, interactionResponse api.InteractionResponse) restclient.RestError {
	compiledRoute, err := restclient.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return restclient.NewError(nil, err)
	}

	body, err := interactionResponse.ToBody()
	if err != nil {
		return restclient.NewError(nil, err)
	}

	return r.Do(compiledRoute, body, nil)
}

// UpdateInteractionResponse used to edit the initial response on an api.Interaction
func (r *restClientImpl) UpdateInteractionResponse(applicationID api.Snowflake, interactionToken string, messageUpdate api.MessageUpdate) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, body, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}
	return
}

// DeleteInteractionResponse used to delete the initial response on an api.Interaction
func (r *restClientImpl) DeleteInteractionResponse(applicationID api.Snowflake, interactionToken string) restclient.RestError {
	compiledRoute, err := restclient.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

// CreateFollowupMessage used to send a followup api.Message to an api.Interaction
func (r *restClientImpl) CreateFollowupMessage(applicationID api.Snowflake, interactionToken string, messageCreate api.MessageCreate) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, body, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}

	return
}

// UpdateFollowupMessage used to edit a followup api.Message from an api.Interaction
func (r *restClientImpl) UpdateFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake, messageUpdate api.MessageUpdate) (message *api.Message, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, body, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, api.CacheStrategyNoWs)
	}

	return
}

// DeleteFollowupMessage used to delete a followup api.Message from an api.Interaction
func (r *restClientImpl) DeleteFollowupMessage(applicationID api.Snowflake, interactionToken string, messageID api.Snowflake) restclient.RestError {
	compiledRoute, err := restclient.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

func (r *restClientImpl) GetGuildTemplate(templateCode string) (guildTemplate *api.GuildTemplate, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildTemplate.Compile(nil, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetGuildTemplates(guildID api.Snowflake) (guildTemplates []*api.GuildTemplate, rErr restclient.RestError) {
	compiledRoute, err := restclient.GetGuildTemplates.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplates)
	if rErr == nil {
		for _, guildTemplate := range guildTemplates {
			guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, api.CacheStrategyNoWs)
		}
	}
	return
}

func (r *restClientImpl) CreateGuildTemplate(guildID api.Snowflake, createGuildTemplate api.CreateGuildTemplate) (guildTemplate *api.GuildTemplate, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateGuildTemplate.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, createGuildTemplate, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate api.CreateGuildFromTemplate) (guild *api.Guild, rErr restclient.RestError) {
	compiledRoute, err := restclient.CreateGuildFromTemplate.Compile(nil, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var fullGuild *api.FullGuild
	rErr = r.Do(compiledRoute, createGuildFromTemplate, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, api.CacheStrategyNoWs)
	}

	return
}

func (r *restClientImpl) SyncGuildTemplate(guildID api.Snowflake, templateCode string) (guildTemplate *api.GuildTemplate, rErr restclient.RestError) {
	compiledRoute, err := restclient.SyncGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) UpdateGuildTemplate(guildID api.Snowflake, templateCode string, updateGuildTemplate api.UpdateGuildTemplate) (guildTemplate *api.GuildTemplate, rErr restclient.RestError) {
	compiledRoute, err := restclient.UpdateGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, updateGuildTemplate, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, api.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) DeleteGuildTemplate(guildID api.Snowflake, templateCode string) (guildTemplate *api.GuildTemplate, rErr restclient.RestError) {
	compiledRoute, err := restclient.DeleteGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, api.CacheStrategyNoWs)
	}
	return
}
