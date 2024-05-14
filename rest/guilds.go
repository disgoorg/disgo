package rest

import (
	"time"

	"github.com/disgoorg/disgo/internal/slicehelper"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ Guilds = (*guildImpl)(nil)

func NewGuilds(client Client) Guilds {
	return &guildImpl{client: client}
}

type Guilds interface {
	GetGuild(guildID snowflake.ID, withCounts bool, opts ...RequestOpt) (*discord.RestGuild, error)
	GetGuildPreview(guildID snowflake.ID, opts ...RequestOpt) (*discord.GuildPreview, error)
	CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (*discord.RestGuild, error)
	UpdateGuild(guildID snowflake.ID, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (*discord.RestGuild, error)
	DeleteGuild(guildID snowflake.ID, opts ...RequestOpt) error

	GetGuildVanityURL(guildID snowflake.ID, opts ...RequestOpt) (*discord.PartialInvite, error)

	CreateGuildChannel(guildID snowflake.ID, guildChannelCreate discord.GuildChannelCreate, opts ...RequestOpt) (discord.GuildChannel, error)
	GetGuildChannels(guildID snowflake.ID, opts ...RequestOpt) ([]discord.GuildChannel, error)
	UpdateChannelPositions(guildID snowflake.ID, guildChannelPositionUpdates []discord.GuildChannelPositionUpdate, opts ...RequestOpt) error

	GetRoles(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Role, error)
	GetRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (*discord.Role, error)
	CreateRole(guildID snowflake.ID, createRole discord.RoleCreate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRole(guildID snowflake.ID, roleID snowflake.ID, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRolePositions(guildID snowflake.ID, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) ([]discord.Role, error)
	DeleteRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error

	GetBans(guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) ([]discord.Ban, error)
	GetBansPage(guildID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[discord.Ban]
	GetBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*discord.Ban, error)
	AddBan(guildID snowflake.ID, userID snowflake.ID, deleteMessageDuration time.Duration, opts ...RequestOpt) error
	DeleteBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	BulkBan(guildID snowflake.ID, ban discord.BulkBan, opts ...RequestOpt) (*discord.BulkBanResult, error)

	GetIntegrations(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Integration, error)
	DeleteIntegration(guildID snowflake.ID, integrationID snowflake.ID, opts ...RequestOpt) error

	GetGuildPruneCount(guildID snowflake.ID, days int, includeRoles []snowflake.ID, opts ...RequestOpt) (*discord.GuildPruneResult, error)
	BeginGuildPrune(guildID snowflake.ID, guildPrune discord.GuildPrune, opts ...RequestOpt) (*discord.GuildPruneResult, error)

	GetAllWebhooks(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Webhook, error)

	GetGuildVoiceRegions(guildID snowflake.ID, opts ...RequestOpt) ([]discord.VoiceRegion, error)

	GetAuditLog(guildID snowflake.ID, userID snowflake.ID, actionType discord.AuditLogEvent, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (*discord.AuditLog, error)
	GetAuditLogPage(guildID snowflake.ID, userID snowflake.ID, actionType discord.AuditLogEvent, startID snowflake.ID, limit int, opts ...RequestOpt) AuditLogPage

	GetGuildWelcomeScreen(guildID snowflake.ID, opts ...RequestOpt) (*discord.GuildWelcomeScreen, error)
	UpdateGuildWelcomeScreen(guildID snowflake.ID, screenUpdate discord.GuildWelcomeScreenUpdate, opts ...RequestOpt) (*discord.GuildWelcomeScreen, error)

	GetGuildOnboarding(guildID snowflake.ID, opts ...RequestOpt) (*discord.GuildOnboarding, error)
	UpdateGuildOnboarding(guildID snowflake.ID, onboardingUpdate discord.GuildOnboardingUpdate, opts ...RequestOpt) (*discord.GuildOnboarding, error)
}

type guildImpl struct {
	client Client
}

func (s *guildImpl) GetGuild(guildID snowflake.ID, withCounts bool, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	values := discord.QueryValues{
		"with_counts": withCounts,
	}
	err = s.client.Do(GetGuild.Compile(values, guildID), nil, &guild, opts...)
	return
}

func (s *guildImpl) GetGuildPreview(guildID snowflake.ID, opts ...RequestOpt) (guildPreview *discord.GuildPreview, err error) {
	err = s.client.Do(GetGuildPreview.Compile(nil, guildID), nil, &guildPreview, opts...)
	return
}

func (s *guildImpl) CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	err = s.client.Do(CreateGuild.Compile(nil), guildCreate, &guild, opts...)
	return
}

func (s *guildImpl) UpdateGuild(guildID snowflake.ID, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	err = s.client.Do(UpdateGuild.Compile(nil, guildID), guildUpdate, &guild, opts...)
	return
}

func (s *guildImpl) DeleteGuild(guildID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGuild.Compile(nil, guildID), nil, nil, opts...)
}

func (s *guildImpl) GetGuildVanityURL(guildID snowflake.ID, opts ...RequestOpt) (partialInvite *discord.PartialInvite, err error) {
	err = s.client.Do(GetGuildVanityURL.Compile(nil, guildID), nil, &partialInvite, opts...)
	return
}

func (s *guildImpl) CreateGuildChannel(guildID snowflake.ID, guildChannelCreate discord.GuildChannelCreate, opts ...RequestOpt) (guildChannel discord.GuildChannel, err error) {
	var ch discord.UnmarshalChannel
	err = s.client.Do(CreateGuildChannel.Compile(nil, guildID), guildChannelCreate, &ch, opts...)
	if err == nil {
		guildChannel = ch.Channel.(discord.GuildChannel)
	}
	return
}

func (s *guildImpl) GetGuildChannels(guildID snowflake.ID, opts ...RequestOpt) (channels []discord.GuildChannel, err error) {
	var chs []discord.UnmarshalChannel
	err = s.client.Do(GetGuildChannels.Compile(nil, guildID), nil, &chs, opts...)
	if err == nil {
		channels = make([]discord.GuildChannel, len(chs))
		for i := range chs {
			channels[i] = chs[i].Channel.(discord.GuildChannel)
		}
	}
	return
}

func (s *guildImpl) UpdateChannelPositions(guildID snowflake.ID, guildChannelPositionUpdates []discord.GuildChannelPositionUpdate, opts ...RequestOpt) error {
	return s.client.Do(UpdateChannelPositions.Compile(nil, guildID), guildChannelPositionUpdates, nil, opts...)
}

func (s *guildImpl) GetRoles(guildID snowflake.ID, opts ...RequestOpt) (roles []discord.Role, err error) {
	err = s.client.Do(GetRoles.Compile(nil, guildID), nil, &roles, opts...)
	if err == nil {
		for i := range roles {
			roles[i].GuildID = guildID
		}
	}
	return
}

func (s *guildImpl) GetRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (role *discord.Role, err error) {
	err = s.client.Do(GetRole.Compile(nil, guildID, roleID), nil, &role, opts...)
	if err == nil {
		role.GuildID = guildID
	}
	return
}

func (s *guildImpl) CreateRole(guildID snowflake.ID, createRole discord.RoleCreate, opts ...RequestOpt) (role *discord.Role, err error) {
	err = s.client.Do(CreateRole.Compile(nil, guildID), createRole, &role, opts...)
	if err == nil {
		role.GuildID = guildID
	}
	return
}

func (s *guildImpl) UpdateRole(guildID snowflake.ID, roleID snowflake.ID, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (role *discord.Role, err error) {
	err = s.client.Do(UpdateRole.Compile(nil, guildID, roleID), roleUpdate, &role, opts...)
	if err == nil {
		role.GuildID = guildID
	}
	return
}

func (s *guildImpl) UpdateRolePositions(guildID snowflake.ID, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) (roles []discord.Role, err error) {
	err = s.client.Do(UpdateRolePositions.Compile(nil, guildID), rolePositionUpdates, &roles, opts...)
	if err == nil {
		for i := range roles {
			roles[i].GuildID = guildID
		}
	}
	return
}

func (s *guildImpl) DeleteRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteRole.Compile(nil, guildID, roleID), nil, nil, opts...)
}

func (s *guildImpl) GetBans(guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (bans []discord.Ban, err error) {
	values := discord.QueryValues{}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetBans.Compile(values, guildID), nil, &bans, opts...)
	return
}

func (s *guildImpl) GetBansPage(guildID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[discord.Ban] {
	return Page[discord.Ban]{
		getItemsFunc: func(before snowflake.ID, after snowflake.ID) (bans []discord.Ban, err error) {
			return s.GetBans(guildID, before, after, limit, opts...)
		},
		getIDFunc: func(ban discord.Ban) snowflake.ID {
			return ban.User.ID
		},
		ID: startID,
	}
}

func (s *guildImpl) GetBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (ban *discord.Ban, err error) {
	err = s.client.Do(GetBan.Compile(nil, guildID, userID), nil, &ban, opts...)
	return
}

func (s *guildImpl) AddBan(guildID snowflake.ID, userID snowflake.ID, deleteMessageDuration time.Duration, opts ...RequestOpt) error {
	return s.client.Do(AddBan.Compile(nil, guildID, userID), discord.AddBan{DeleteMessageSeconds: int(deleteMessageDuration.Seconds())}, nil, opts...)
}

func (s *guildImpl) DeleteBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteBan.Compile(nil, guildID, userID), nil, nil, opts...)
}

func (s *guildImpl) BulkBan(guildID snowflake.ID, ban discord.BulkBan, opts ...RequestOpt) (result *discord.BulkBanResult, err error) {
	err = s.client.Do(BulkBan.Compile(nil, guildID), ban, &result, opts...)
	return
}

func (s *guildImpl) GetIntegrations(guildID snowflake.ID, opts ...RequestOpt) (integrations []discord.Integration, err error) {
	err = s.client.Do(GetIntegrations.Compile(nil, guildID), nil, &integrations, opts...)
	return
}

func (s *guildImpl) DeleteIntegration(guildID snowflake.ID, integrationID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteIntegration.Compile(nil, guildID, integrationID), nil, nil, opts...)
}

func (s *guildImpl) GetGuildPruneCount(guildID snowflake.ID, days int, includeRoles []snowflake.ID, opts ...RequestOpt) (result *discord.GuildPruneResult, err error) {
	values := discord.QueryValues{
		"days":          days,
		"include_roles": slicehelper.JoinSnowflakes(includeRoles),
	}
	err = s.client.Do(GetGuildPruneCount.Compile(values, guildID), nil, &result, opts...)
	return
}

func (s *guildImpl) BeginGuildPrune(guildID snowflake.ID, guildPrune discord.GuildPrune, opts ...RequestOpt) (result *discord.GuildPruneResult, err error) {
	err = s.client.Do(BeginGuildPrune.Compile(nil, guildID), guildPrune, &result, opts...)
	return
}

func (s *guildImpl) GetAllWebhooks(guildID snowflake.ID, opts ...RequestOpt) (webhooks []discord.Webhook, err error) {
	var whs []discord.UnmarshalWebhook
	err = s.client.Do(GetGuildWebhooks.Compile(nil, guildID), nil, &whs, opts...)
	if err == nil {
		webhooks = make([]discord.Webhook, len(whs))
		for i := range whs {
			webhooks[i] = whs[i].Webhook
		}
	}
	return
}

func (s *guildImpl) GetGuildVoiceRegions(guildID snowflake.ID, opts ...RequestOpt) (regions []discord.VoiceRegion, err error) {
	err = s.client.Do(GetGuildVoiceRegions.Compile(nil, guildID), nil, &regions, opts...)
	return
}

func (s *guildImpl) GetAuditLog(guildID snowflake.ID, userID snowflake.ID, actionType discord.AuditLogEvent, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (auditLog *discord.AuditLog, err error) {
	values := discord.QueryValues{}
	if userID != 0 {
		values["user_id"] = userID
	}
	if actionType != 0 {
		values["action_type"] = actionType
	}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetAuditLogs.Compile(values, guildID), nil, &auditLog, opts...)
	return
}

func (s *guildImpl) GetAuditLogPage(guildID snowflake.ID, userID snowflake.ID, actionType discord.AuditLogEvent, startID snowflake.ID, limit int, opts ...RequestOpt) AuditLogPage {
	return AuditLogPage{
		getItems: func(before snowflake.ID, after snowflake.ID) (discord.AuditLog, error) {
			log, err := s.GetAuditLog(guildID, userID, actionType, before, after, limit, opts...)
			var finalLog discord.AuditLog
			if log != nil {
				finalLog = *log
			}
			return finalLog, err
		},
		ID: startID,
	}
}

func (s *guildImpl) GetGuildWelcomeScreen(guildID snowflake.ID, opts ...RequestOpt) (welcomeScreen *discord.GuildWelcomeScreen, err error) {
	err = s.client.Do(GetGuildWelcomeScreen.Compile(nil, guildID), nil, &welcomeScreen, opts...)
	return
}

func (s *guildImpl) UpdateGuildWelcomeScreen(guildID snowflake.ID, screenUpdate discord.GuildWelcomeScreenUpdate, opts ...RequestOpt) (welcomeScreen *discord.GuildWelcomeScreen, err error) {
	err = s.client.Do(UpdateGuildWelcomeScreen.Compile(nil, guildID), screenUpdate, &welcomeScreen, opts...)
	return
}

func (s *guildImpl) GetGuildOnboarding(guildID snowflake.ID, opts ...RequestOpt) (onboarding *discord.GuildOnboarding, err error) {
	err = s.client.Do(GetGuildOnboarding.Compile(nil, guildID), nil, &onboarding, opts...)
	return
}

func (s *guildImpl) UpdateGuildOnboarding(guildID snowflake.ID, onboardingUpdate discord.GuildOnboardingUpdate, opts ...RequestOpt) (guildOnboarding *discord.GuildOnboarding, err error) {
	err = s.client.Do(UpdateGuildOnboarding.Compile(nil, guildID), onboardingUpdate, &guildOnboarding, opts...)
	return
}
