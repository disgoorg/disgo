package rest

import (
	"github.com/DisgoOrg/log"
)

var _ Services = (*ServicesImpl)(nil)

func NewServices(logger log.Logger, restClient Client) Services {
	return &ServicesImpl{
		logger:               logger,
		restClient:           restClient,
		applicationService:   NewApplicationService(restClient),
		oauth2Service:        NewOAuth2Service(restClient),
		auditLogService:      NewAuditLogService(restClient),
		gatewayService:       NewGatewayService(restClient),
		guildService:         NewGuildService(restClient),
		channelService:       NewChannelService(restClient),
		interactionService:   NewInteractionService(restClient),
		inviteService:        NewInviteService(restClient),
		guildTemplateService: NewGuildTemplateService(restClient),
		userService:          NewUserService(restClient),
		voiceService:         NewVoiceService(restClient),
		webhookService:       NewWebhookService(restClient),
		stageInstanceService: NewStageInstanceService(restClient),
		emojiService:         NewEmojiService(restClient),
	}
}

// Services is a manager for all of disgo's HTTP requests
type Services interface {
	Close()
	Logger() log.Logger
	RestClient() Client
	ApplicationService() ApplicationService
	OAuth2Service() OAuth2Service
	AuditLogService() AuditLogService
	GatewayService() GatewayService
	GuildService() GuildService
	ChannelService() ChannelService
	InteractionService() InteractionService
	InviteService() InviteService
	GuildTemplateService() GuildTemplateService
	UserService() UserService
	VoiceService() VoiceService
	WebhookService() WebhookService
	StageInstanceService() StageInstanceService
	EmojiService() EmojiService
}

type ServicesImpl struct {
	logger     log.Logger
	restClient Client

	applicationService   ApplicationService
	oauth2Service        OAuth2Service
	auditLogService      AuditLogService
	gatewayService       GatewayService
	guildService         GuildService
	channelService       ChannelService
	interactionService   InteractionService
	inviteService        InviteService
	guildTemplateService GuildTemplateService
	userService          UserService
	voiceService         VoiceService
	webhookService       WebhookService
	stageInstanceService StageInstanceService
	emojiService         EmojiService
}

func (s *ServicesImpl) Close() {
	s.restClient.Close()
}

func (s *ServicesImpl) Logger() log.Logger {
	return s.logger
}

func (s *ServicesImpl) RestClient() Client {
	return s.restClient
}

func (s *ServicesImpl) ApplicationService() ApplicationService {
	return s.applicationService
}

func (s *ServicesImpl) OAuth2Service() OAuth2Service {
	return s.oauth2Service
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

func (s *ServicesImpl) ChannelService() ChannelService {
	return s.channelService
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

func (s *ServicesImpl) EmojiService() EmojiService {
	return s.emojiService
}

type Service interface {
	RestClient() Client
}
