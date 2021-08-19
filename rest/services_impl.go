package rest

import "github.com/DisgoOrg/log"

var _ Services = (*ServicesImpl)(nil)

func NewServices(logger log.Logger, httpClient Client) Services {
	return &ServicesImpl{
		logger:               logger,
		httpClient:           httpClient,
		applicationService:   nil,
		auditLogService:      nil,
		gatewayService:       nil,
		guildService:         nil,
		channelsService:      nil,
		interactionService:   nil,
		inviteService:        nil,
		guildTemplateService: nil,
		userService:          nil,
		voiceService:         nil,
		webhookService:       nil,
		stageInstanceService: nil,
	}
}

type ServicesImpl struct {
	logger     log.Logger
	httpClient Client

	applicationService   ApplicationService
	auditLogService      AuditLogService
	gatewayService       GatewayService
	guildService         GuildService
	channelsService      ChannelsService
	interactionService   InteractionService
	inviteService        InviteService
	guildTemplateService GuildTemplateService
	userService          UserService
	voiceService         VoiceService
	webhookService       WebhookService
	stageInstanceService StageInstanceService
}

func (s *ServicesImpl) Close() {
	s.httpClient.Close()
}

func (s *ServicesImpl) Logger() log.Logger {
	return s.logger
}

func (s *ServicesImpl) HTTPClient() Client {
	return s.httpClient
}

func (s *ServicesImpl) ApplicationService() ApplicationService {
	return s.applicationService
}
func (s *ServicesImpl) AuditLogService() AuditLogService {
	return s.auditLogService
}
func (s *ServicesImpl) GatewayService() GatewayService {
	return s.gatewayService
}
func (s *ServicesImpl) GuildService() GuildService {
	return s.guildService
}
func (s *ServicesImpl) ChannelsService() ChannelsService {
	return s.channelsService
}
func (s *ServicesImpl) InteractionService() InteractionService {
	return s.interactionService
}
func (s *ServicesImpl) InviteService() InviteService {
	return s.inviteService
}
func (s *ServicesImpl) GuildTemplateService() GuildTemplateService {
	return s.guildTemplateService
}
func (s *ServicesImpl) UserService() UserService {
	return s.userService
}
func (s *ServicesImpl) VoiceService() VoiceService {
	return s.voiceService
}
func (s *ServicesImpl) WebhookService() WebhookService {
	return s.webhookService
}
func (s *ServicesImpl) StageInstanceService() StageInstanceService {
	return s.stageInstanceService
}

/*
func NewRestClient(disgo core.Disgo, httpClient *httpserver.Client) Services {
	if httpClient == nil {
		httpClient = httpserver.DefaultClient
	}
	return &restClientImpl{
		Services: restclient.NewRestClient(httpClient, disgo.Logger(), UserAgent, httpserver.Header{"Authorization": []string{"Bot " + disgo.Token()}}),
		disgo:      disgo,
	}
}

 restClientImpl is the rest client implementation used for HTTP requests to discord
type restClientImpl struct {
	restclient.Services
	disgo core.Disgo
}

func (r *restClientImpl) ApplicationService() ApplicationService {
	panic("implement me")
}

func (r *restClientImpl) AuditLogService() AuditLogService {
	panic("implement me")
}

func (r *restClientImpl) EmojiService() EmojiService {
	panic("implement me")
}

func (r *restClientImpl) GatewayService() GatewayService {
	panic("implement me")
}

func (r *restClientImpl) GuildService() GuildService {
	panic("implement me")
}

func (r *restClientImpl) ChannelsService() ChannelsService {
	panic("implement me")
}

func (r *restClientImpl) InteractionService() InteractionService {
	panic("implement me")
}

func (r *restClientImpl) InviteService() InviteService {
	panic("implement me")
}

func (r *restClientImpl) GuildTemplateService() GuildTemplateService {
	panic("implement me")
}

func (r *restClientImpl) UserService() UserService {
	panic("implement me")
}

func (r *restClientImpl) VoiceService() VoiceService {
	panic("implement me")
}

func (r *restClientImpl) WebhookService() WebhookService {
	panic("implement me")
}

 Disgo returns the api.Disgo instance
func (r *restClientImpl) Disgo() core.Disgo {
	return r.disgo
}

 Close cleans up the httpserver managers connections
func (r *restClientImpl) Close() {
	r.Client().CloseIdleConnections()
}

 DoWithHeaders executes a rest request with custom headers
func (r *restClientImpl) DoWithHeaders(route *restclient.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader httpserver.Header) (rErr rest.Error) {
	err := r.Services.DoWithHeaders(route, rqBody, rsBody, customHeader)
	rErr = restclient.NewError(nil, err)
	 TODO reimplement events.HTTPRequestEvent
	/*r.Disgo().EventManager().Dispatch(&events.HTTPRequestEvent{
		GenericEvent: events.NewEvent(r.Disgo(), 0),
		Request:      rq,
		Response:     rs,
	})

	 TODO reimplement api.ErrorResponse unmarshalling
	/*
		var errorRs api.ErrorResponse
				if err = json.Unmarshal(rawRsBody, &errorRs); err != nil {
					r.Disgo().Logger().Errorf("rest.Error unmarshalling rest.Error response. code: %d, rest.Error: %s", rs.StatusCode, err)
					return err
				}
				return fmt.Errorf("request to %s failed. statuscode: %d, errorcode: %d, message_events: %s", rq.URL, rs.StatusCode, errorRs.Code, errorRs.Message)

	return
}

 GetUser fetches the specific user
func (r *restClientImpl) GetUser(userID discord.Snowflake) (user *discord.User, rErr rest.Error) {
	compiledRoute, err := restclient.GetUser.Compile(nil, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &user)
	if rErr == nil {
		user = r.Disgo().EntityBuilder().CreateUser(user, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetSelfUser() (selfUser *discord.SelfUser, rErr rest.Error) {
	compiledRoute, err := restclient.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	var user *discord.User
	rErr = r.Do(compiledRoute, nil, &user)
	if rErr == nil {
		selfUser = &discord.SelfUser{User: r.Disgo().EntityBuilder().CreateUser(user, core.CacheStrategyNoWs)}
	}
	return
}

func (r *restClientImpl) UpdateSelfUser(updateSelfUser discord.UpdateSelfUser) (selfUser *discord.SelfUser, rErr rest.Error) {
	compiledRoute, err := restclient.GetSelfUser.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	var user *discord.User
	rErr = r.Do(compiledRoute, updateSelfUser, &user)
	if rErr == nil {
		selfUser = &discord.SelfUser{User: r.Disgo().EntityBuilder().CreateUser(user, core.CacheStrategyNoWs)}
	}
	return
}

func (r *restClientImpl) GetGuilds(before int, after int, limit int) (guilds []*discord.PartialGuild, rErr rest.Error) {
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

func (r *restClientImpl) LeaveGuild(guildID discord.Snowflake) rest.Error {
	compiledRoute, err := restclient.LeaveGuild.Compile(nil, guildID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

func (r *restClientImpl) GetDMChannels() (dmChannels []*discord.DMChannel, rErr rest.Error) {
	compiledRoute, err := restclient.GetDMChannels.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var channels []*api.Channel
	rErr = r.Do(compiledRoute, nil, &channels)
	if rErr == nil {
		dmChannels = make([]*api.DMChannel, len(channels))
		for i, channel := range channels {
			dmChannels[i] = r.Disgo().EntityBuilder().CreateDMChannel(channel, core.CacheStrategyNoWs)
		}
	}
	return
}

 CreateDMChannel opens a new api.DMChannel to an api.User
func (r *restClientImpl) CreateDMChannel(userID discord.Snowflake) (channel *discord.DMChannel, rErr rest.Error) {
	compiledRoute, err := restclient.CreateDMChannel.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, api.CreateDMChannel{RecipientID: userID}, &channel)
	if rErr == nil {
		channel = r.Disgo().EntityBuilder().CreateDMChannel(&channel.MessageChannel.Channel, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetMessage(channelID discord.Snowflake, messageID discord.Snowflake) (message *entities.Message, rErr rest.Error) {
	compiledRoute, err := restclient.GetMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}
	return
}

 CreateMessage lets you send an api.Message to an api.MessageChannel
func (r *restClientImpl) CreateMessage(channelID discord.Snowflake, messageCreate api.MessageCreate) (message *entities.Message, rErr rest.Error) {
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
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}
	return
}

 UpdateMessage lets you edit an api.Message
func (r *restClientImpl) UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate api.MessageUpdate) (message *entities.Message, rErr rest.Error) {
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
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}
	return
}

 DeleteMessage lets you delete an api.Message
func (r *restClientImpl) DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().UncacheMessage(channelID, messageID)
	}
	return
}

 BulkDeleteMessages lets you bulk delete api.Message(s)
func (r *restClientImpl) BulkDeleteMessages(channelID discord.Snowflake, messageIDs ...discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.BulkDeleteMessage.Compile(nil, channelID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, entities.MessageBulkDelete{Messages: messageIDs}, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		 TODO: check here if no err means all messages deleted
		for _, messageID := range messageIDs {
			r.Disgo().Cache().UncacheMessage(channelID, messageID)
		}
	}
	return
}

 CrosspostMessage lets you crosspost an api.Message in a channel with type api.ChannelTypeNews
func (r *restClientImpl) CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake) (message *entities.Message, rErr rest.Error) {
	compiledRoute, err := restclient.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &message)
	if rErr == nil {
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetGuild(guildID discord.Snowflake, withCounts bool) (guild *api.Guild, rErr rest.Error) {
	var queryParams = restclient.QueryValues{
		"with_counts": withCounts,
	}
	compiledRoute, err := restclient.GetGuild.Compile(queryParams, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, restclient.NewError(nil, err))
	}

	var fullGuild *discord.FullGuild
	rErr = r.Do(compiledRoute, nil, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetGuildPreview(guildID discord.Snowflake) (guildPreview *api.GuildPreview, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildPreview.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, restclient.NewError(nil, err))
	}

	rErr = r.Do(compiledRoute, nil, &guildPreview)
	return
}

func (r *restClientImpl) CreateGuild(createGuild discord.CreateGuild) (guild *discord.Guild, rErr rest.Error) {
	compiledRoute, err := restclient.CreateGuild.Compile(nil)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var fullGuild *discord.FullGuild
	rErr = r.Do(compiledRoute, createGuild, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) UpdateGuild(guildID discord.Snowflake, updateGuild api.UpdateGuild) (guild *api.Guild, rErr rest.Error) {
	compiledRoute, err := restclient.CreateGuild.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var fullGuild *discord.FullGuild
	rErr = r.Do(compiledRoute, updateGuild, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) DeleteGuild(guildID discord.Snowflake) rest.Error {
	compiledRoute, err := restclient.DeleteGuild.Compile(nil, guildID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

 GetMember fetches the specific member
func (r *restClientImpl) GetMember(guildID discord.Snowflake, userID discord.Snowflake) (member *entities.Member, rErr rest.Error) {
	compiledRoute, err := restclient.GetMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, core.CacheStrategyNoWs)
	}
	return
}

 GetMembers fetches all members for a guild
func (r *restClientImpl) GetMembers(guildID discord.Snowflake) (members []*discord.Member, rErr rest.Error) {
	compiledRoute, err := restclient.GetMembers.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &members)
	if rErr == nil {
		for _, member := range members {
			member = r.Disgo().EntityBuilder().CreateMember(guildID, member, core.CacheStrategyNoWs)
		}
	}
	return
}

func (r *restClientImpl) SearchMembers(guildID discord.Snowflake, query string, limit int) (members []*discord.Member, rErr rest.Error) {
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
		members = make([]*discord.Member, len(members))
		for i, member := range members {
			members[i] = r.Disgo().EntityBuilder().CreateMember(guildID, member, core.CacheStrategyNoWs)
		}
	}
	return
}

 AddMember adds a member to the guild with the oauth2 access BotToken. requires api.PermissionCreateInstantInvite
func (r *restClientImpl) AddMember(guildID discord.Snowflake, userID discord.Snowflake, addMember discord.AddMember) (member *discord.Member, rErr rest.Error) {
	compiledRoute, err := restclient.AddMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, addMember, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, core.CacheStrategyNoWs)
	}
	return
}

 RemoveMember kicks an api.Member from the api.Guild. requires api.PermissionKickMembers
func (r *restClientImpl) RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, reason string) (rErr rest.Error) {
	var params restclient.QueryValues
	if reason != "" {
		params = restclient.QueryValues{"reason": reason}
	}
	compiledRoute, err := restclient.RemoveMember.Compile(params, guildID, userID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().Uncache(guildID, userID)
	}
	return
}

 UpdateMember updates an api.Member
func (r *restClientImpl) UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, updateMember entities.UpdateMember) (member *entities.Member, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, updateMember, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, core.CacheStrategyNoWs)
	}
	return
}

 UpdateSelfNick updates the bots nickname in a guild
func (r *restClientImpl) UpdateSelfNick(guildID discord.Snowflake, nick string) (newNick *string, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateSelfNick.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	var updateNick *discord.UpdateSelfNick
	rErr = r.Do(compiledRoute, &discord.UpdateSelfNick{Nick: nick}, &updateNick)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		var nick *string
		if updateNick.Nick == "" {
			nick = nil
		}
		r.Disgo().Cache().Member(guildID, r.Disgo().ClientID()).Nick = nick
		newNick = nick
	}
	return
}

 MoveMember moves/kicks the api.Member to/from an api.VoiceChannel
func (r *restClientImpl) MoveMember(guildID discord.Snowflake, userID discord.Snowflake, channelID *discord.Snowflake) (member *entities.Member, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateMember.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, discord.MoveMember{ChannelID: channelID}, &member)
	if rErr == nil {
		member = r.Disgo().EntityBuilder().CreateMember(guildID, member, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetAuditLog(guildID discord.Snowflake, userID discord.Snowflake, actionType discord.AuditLogEvent, before discord.Snowflake, limit int) (auditLog *discord.AuditLog, rErr rest.Error) {
	values := restclient.QueryValues{}
	if guildID != "" {
		values["guild_id"] = guildID
	}
	if userID != "" {
		values["user_id"] = userID
	}
	if actionType != 0 {
		values["action_type"] = actionType
	}
	if before != "" {
		values["before"] = guildID
	}
	if limit != 0 {
		values["limit"] = limit
	}
	compiledRoute, err := restclient.GetAuditLogs.Compile(values, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &auditLog)
	if rErr == nil {
		auditLog = r.Disgo().EntityBuilder().CreateAuditLog(guildID, discord.AuditLogFilterOptions{
			UserID:     userID,
			ActionType: actionType,
			Before:     before,
			Limit:      limit,
		}, auditLog, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetIntegrations(guildID discord.Snowflake) (integrations []*discord.Integration, rErr rest.Error) {
	compiledRoute, err := restclient.GetIntegrations.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &integrations)
	if rErr == nil {
		for _, integration := range integrations {
			integration = r.Disgo().EntityBuilder().CreateIntegration(guildID, integration, core.CacheStrategyNoWs)
		}
	}
	return
}

func (r *restClientImpl) DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake) rest.Error {
	compiledRoute, err := restclient.DeleteIntegration.Compile(nil, guildID, integrationID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

func (r *restClientImpl) GetBans(guildID discord.Snowflake) (bans []discord.Ban, rErr rest.Error) {
	compiledRoute, err := restclient.GetBans.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &bans)
	if rErr == nil {
		for _, ban := range bans {
			ban.User = r.Disgo().EntityBuilder().CreateUser(ban.User, core.CacheStrategyNoWs)
		}
	}
	return
}
func (r *restClientImpl) GetBan(guildID discord.Snowflake, userID discord.Snowflake) (ban *discord.Ban, rErr rest.Error) {
	compiledRoute, err := restclient.GetBan.Compile(nil, guildID, userID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &ban)
	if rErr == nil {
		ban.User = r.Disgo().EntityBuilder().CreateUser(ban.User, core.CacheStrategyNoWs)
	}
	return
}
func (r *restClientImpl) AddBan(guildID discord.Snowflake, userID discord.Snowflake, reason string, deleteMessageDays int) rest.Error {
	compiledRoute, err := restclient.AddBan.Compile(nil, guildID, userID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, discord.AddBan{DeleteMessageDays: deleteMessageDays, Reason: reason}, nil)
}

func (r *restClientImpl) DeleteBan(guildID discord.Snowflake, userID discord.Snowflake) rest.Error {
	compiledRoute, err := restclient.DeleteBan.Compile(nil, guildID, userID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

 AddMemberRole adds an api.Role to an api.Member
func (r *restClientImpl) AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.AddMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		member := r.Disgo().Cache().Member(guildID, userID)
		if member != nil {
			member.RoleIDs = append(member.RoleIDs, roleID)
		}
	}
	return
}

 RemoveMemberRole removes an api.Role(s) from an api.Member
func (r *restClientImpl) RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.RemoveMemberRole.Compile(nil, guildID, userID, roleID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
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

 GetRoles fetches all api.Role(s) from an api.Guild
func (r *restClientImpl) GetRoles(guildID discord.Snowflake) (roles []*entities.Role, rErr rest.Error) {
	compiledRoute, err := restclient.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &roles)
	if rErr == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, core.CacheStrategyNoWs)
		}
	}
	return
}

 CreateRole creates a new role for a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) CreateRole(guildID discord.Snowflake, createRole discord.CreateRole) (newRole *discord.Role, rErr rest.Error) {
	compiledRoute, err := restclient.CreateRole.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, createRole, &newRole)
	if rErr == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, core.CacheStrategyNoWs)
	}
	return
}

 UpdateRole updates a role from a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, role discord.UpdateRole) (newRole *discord.Role, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, role, &newRole)
	if rErr == nil {
		newRole = r.Disgo().EntityBuilder().CreateRole(guildID, newRole, core.CacheStrategyNoWs)
	}
	return
}

 UpdateRolePositions updates the position of a role from a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) UpdateRolePositions(guildID discord.Snowflake, roleUpdates ...discord.UpdateRolePosition) (roles []*entities.Role, rErr rest.Error) {
	compiledRoute, err := restclient.GetRoles.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, roleUpdates, &roles)
	if rErr == nil {
		for _, role := range roles {
			role = r.Disgo().EntityBuilder().CreateRole(guildID, role, core.CacheStrategyNoWs)
		}
	}
	return
}

 DeleteRole deletes a role from a guild. Requires api.PermissionManageRoles
func (r *restClientImpl) DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		r.disgo.Cache().Uncache(guildID, roleID)
	}
	return
}

 AddReaction lets you add a reaction to an api.Message
func (r *restClientImpl) AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string) rest.Error {
	compiledRoute, err := restclient.AddReaction.Compile(nil, channelID, messageID, internal.normalizeEmoji(emoji))
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

 RemoveOwnReaction lets you remove your own reaction from an api.Message
func (r *restClientImpl) RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string) rest.Error {
	compiledRoute, err := restclient.RemoveOwnReaction.Compile(nil, channelID, messageID, internal.normalizeEmoji(emoji))
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

 RemoveUserReaction lets you remove a specific reaction from an api.User from an api.Message
func (r *restClientImpl) RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake) rest.Error {
	compiledRoute, err := restclient.RemoveUserReaction.Compile(nil, channelID, messageID, internal.normalizeEmoji(emoji), userID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

 GetGlobalCommands gets you all global api.ApplicationCommand(s)
func (r *restClientImpl) GetGlobalCommands(applicationID discord.Snowflake) (commands []*discord.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.GetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &commands)
	if rErr == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, core.CacheStrategyNoWs)
		}
	}
	return
}

 GetGlobalCommand gets you a specific global global api.ApplicationCommand
func (r *restClientImpl) GetGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake) (cmd *discord.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.GetGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, core.CacheStrategyNoWs)
	}
	return
}

 CreateGlobalCommand lets you create a new global api.ApplicationCommand
func (r *restClientImpl) CreateGlobalCommand(applicationID discord.Snowflake, command discord.CommandCreate) (cmd *discord.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.CreateGlobalCommand.Compile(nil, applicationID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, core.CacheStrategyNoWs)
	}
	return
}

 SetGlobalCommands lets you override all global api.ApplicationCommand
func (r *restClientImpl) SetGlobalCommands(applicationID discord.Snowflake, commands ...discord.CommandCreate) (cmds []*discord.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.SetGlobalCommands.Compile(nil, applicationID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	if len(commands) > 100 {
		err = ErrMaxCommands
		return
	}
	rErr = r.Do(compiledRoute, commands, &cmds)
	if rErr == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, core.CacheStrategyNoWs)
		}
	}
	return
}

 UpdateGlobalCommand lets you edit a specific global api.ApplicationCommand
func (r *restClientImpl) UpdateGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, command api.CommandUpdate) (cmd *api.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGlobalCommand(cmd, core.CacheStrategyNoWs)
	}
	return
}

 DeleteGlobalCommand lets you delete a specific global api.ApplicationCommand
func (r *restClientImpl) DeleteGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.DeleteGlobalCommand.Compile(nil, applicationID, commandID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().Uncache(commandID)
	}
	return
}

 GetGuildCommands gets you all api.ApplicationCommand(s) from an api.Guild
func (r *restClientImpl) GetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake) (commands []*api.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &commands)
	if rErr == nil {
		for _, cmd := range commands {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, core.CacheStrategyNoWs)
		}
	}
	return
}

 CreateGuildCommand lets you create a new api.ApplicationCommand in an api.Guild
func (r *restClientImpl) CreateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, command api.CommandCreate) (cmd *api.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.CreateGuildCommand.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, core.CacheStrategyNoWs)
	}
	return
}

 SetGuildCommands lets you override all api.ApplicationCommand(s) in an api.Guild
func (r *restClientImpl) SetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, commands ...api.CommandCreate) (cmds []*api.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.SetGuildCommands.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	if len(commands) > 100 {
		err = ErrMaxCommands
		return
	}
	rErr = r.Do(compiledRoute, commands, &cmds)
	if rErr == nil {
		for _, cmd := range cmds {
			cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, core.CacheStrategyNoWs)
		}
	}
	return
}

 GetGuildCommand gets you a specific api.ApplicationCommand in an api.Guild
func (r *restClientImpl) GetGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (cmd *api.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, core.CacheStrategyNoWs)
	}
	return
}

 UpdateGuildCommand lets you edit a specific api.ApplicationCommand in an api.Guild
func (r *restClientImpl) UpdateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, command api.CommandUpdate) (cmd *api.ApplicationCommand, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, command, &cmd)
	if rErr == nil {
		cmd = r.Disgo().EntityBuilder().CreateGuildCommand(guildID, cmd, core.CacheStrategyNoWs)
	}
	return
}

 DeleteGuildCommand lets you delete a specific api.ApplicationCommand in an api.Guild
func (r *restClientImpl) DeleteGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (rErr rest.Error) {
	compiledRoute, err := restclient.DeleteGuildCommand.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, nil)
	if rErr == nil && core.CacheStrategyNoWs(r.Disgo()) {
		r.Disgo().Cache().Uncache(commandID)
	}
	return
}

 GetGuildCommandsPermissions returns the api.CommandPermission for a all api.ApplicationCommand(s) in an api.Guild
func (r *restClientImpl) GetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake) (cmdsPerms []*api.GuildCommandPermissions, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildCommandPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmdsPerms)
	if rErr == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, core.CacheStrategyNoWs)
		}
	}
	return
}

 GetGuildCommandPermissions returns the api.CommandPermission for a specific api.ApplicationCommand in an api.Guild
func (r *restClientImpl) GetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (cmdPerms *api.GuildCommandPermissions, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildCommandPermission.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, nil, &cmdPerms)
	if rErr == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, core.CacheStrategyNoWs)
	}
	return
}

 SetGuildCommandsPermissions sets the api.GuildCommandPermissions for a all api.ApplicationCommand(s)
func (r *restClientImpl) SetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandsPermissions ...api.SetGuildCommandPermissions) (cmdsPerms []*api.GuildCommandPermissions, rErr rest.Error) {
	compiledRoute, err := restclient.SetGuildCommandsPermissions.Compile(nil, applicationID, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, api.SetGuildCommandsPermissions(commandsPermissions), &cmdsPerms)
	if rErr == nil {
		for _, cmdPerms := range cmdsPerms {
			cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, core.CacheStrategyNoWs)
		}
	}
	return
}

 SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific api.ApplicationCommand
func (r *restClientImpl) SetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions discord.SetGuildCommandPermissions) (cmdPerms *api.GuildCommandPermissions, rErr rest.Error) {
	compiledRoute, err := restclient.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}
	rErr = r.Do(compiledRoute, commandPermissions, &cmdPerms)
	if rErr == nil {
		cmdPerms = r.Disgo().EntityBuilder().CreateGuildCommandPermissions(cmdPerms, core.CacheStrategyNoWs)
	}
	return
}

 SendInteractionResponse used to send the initial response on an api.Interaction
func (r *restClientImpl) SendInteractionResponse(interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse) rest.Error {
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

func (r *restClientImpl) UpdateInteractionResponse(applicationID discord.Snowflake, interactionToken string, messageUpdate api.MessageUpdate) (message *entities.Message, rErr rest.Error) {
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
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) DeleteInteractionResponse(applicationID discord.Snowflake, interactionToken string) rest.Error {
	compiledRoute, err := restclient.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

func (r *restClientImpl) SendFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageCreate api.MessageCreate) (message *entities.Message, rErr rest.Error) {
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
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}

	return
}

func (r *restClientImpl) UpdateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate api.MessageUpdate) (message *entities.Message, rErr rest.Error) {
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
		message = r.Disgo().EntityBuilder().CreateMessage(message, core.CacheStrategyNoWs)
	}

	return
}

func (r *restClientImpl) DeleteFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake) rest.Error {
	compiledRoute, err := restclient.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return restclient.NewError(nil, err)
	}
	return r.Do(compiledRoute, nil, nil)
}

func (r *restClientImpl) GetGuildTemplate(templateCode string) (guildTemplate *api.GuildTemplate, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildTemplate.Compile(nil, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) GetGuildTemplates(guildID discord.Snowflake) (guildTemplates []*api.GuildTemplate, rErr rest.Error) {
	compiledRoute, err := restclient.GetGuildTemplates.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplates)
	if rErr == nil {
		for _, guildTemplate := range guildTemplates {
			guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, core.CacheStrategyNoWs)
		}
	}
	return
}

func (r *restClientImpl) CreateGuildTemplate(guildID discord.Snowflake, createGuildTemplate api.CreateGuildTemplate) (guildTemplate *api.GuildTemplate, rErr rest.Error) {
	compiledRoute, err := restclient.CreateGuildTemplate.Compile(nil, guildID)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, createGuildTemplate, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) CreateGuildFromTemplate(templateCode string, createGuildFromTemplate api.CreateGuildFromTemplate) (guild *api.Guild, rErr rest.Error) {
	compiledRoute, err := restclient.CreateGuildFromTemplate.Compile(nil, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	var fullGuild *api.FullGuild
	rErr = r.Do(compiledRoute, createGuildFromTemplate, &fullGuild)
	if rErr == nil {
		guild = r.Disgo().EntityBuilder().CreateGuild(fullGuild, core.CacheStrategyNoWs)
	}

	return
}

func (r *restClientImpl) SyncGuildTemplate(guildID discord.Snowflake, templateCode string) (guildTemplate *api.GuildTemplate, rErr rest.Error) {
	compiledRoute, err := restclient.SyncGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, updateGuildTemplate api.UpdateGuildTemplate) (guildTemplate *api.GuildTemplate, rErr rest.Error) {
	compiledRoute, err := restclient.UpdateGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, updateGuildTemplate, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, core.CacheStrategyNoWs)
	}
	return
}

func (r *restClientImpl) DeleteGuildTemplate(guildID discord.Snowflake, templateCode string) (guildTemplate *api.GuildTemplate, rErr rest.Error) {
	compiledRoute, err := restclient.DeleteGuildTemplate.Compile(nil, guildID, templateCode)
	if err != nil {
		return nil, restclient.NewError(nil, err)
	}

	rErr = r.Do(compiledRoute, nil, &guildTemplate)
	if rErr == nil {
		guildTemplate = r.Disgo().EntityBuilder().CreateGuildTemplate(guildTemplate, core.CacheStrategyNoWs)
	}
	return
}
*/
