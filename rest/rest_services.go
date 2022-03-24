package rest

import (
	"context"
	"net/http"
)

var _ Rest = (*servicesImpl)(nil)

// NewRest returns a new default Rest
func NewRest(restClient Client) Rest {
	return &servicesImpl{
		restClient:                 restClient,
		applicationService:         NewApplicationService(restClient),
		oauth2Service:              NewOAuth2Service(restClient),
		auditLogService:            NewAuditLogService(restClient),
		gatewayService:             NewGatewayService(restClient),
		guildService:               NewGuildService(restClient),
		memberService:              NewMemberService(restClient),
		channelService:             NewChannelService(restClient),
		threadService:              NewThreadService(restClient),
		interactionService:         NewInteractionService(restClient),
		inviteService:              NewInviteService(restClient),
		guildTemplateService:       NewGuildTemplateService(restClient),
		userService:                NewUserService(restClient),
		voiceService:               NewVoiceService(restClient),
		webhookService:             NewWebhookService(restClient),
		stageInstanceService:       NewStageInstanceService(restClient),
		emojiService:               NewEmojiService(restClient),
		stickerService:             NewStickerService(restClient),
		guildScheduledEventService: NewGuildScheduledEventService(restClient),
	}
}

// Rest is a manager for all of disgo's HTTP requests
type Rest interface {
	RestClient() Client
	HTTPClient() *http.Client
	Close(ctx context.Context)

	Application() ApplicationService
	OAuth2() OAuth2Service
	AuditLog() AuditLogService
	Gateway() GatewayService
	Guild() GuildService
	Member() MemberService
	Channel() ChannelService
	Thread() ThreadService
	Interaction() InteractionService
	Invite() InviteService
	GuildTemplate() GuildTemplateService
	User() UserService
	Voice() VoiceService
	Webhook() WebhookService
	StageInstance() StageInstanceService
	Emoji() EmojiService
	Sticker() StickerService
	GuildScheduledEvent() GuildScheduledEventService
}

type servicesImpl struct {
	restClient Client

	applicationService         ApplicationService
	oauth2Service              OAuth2Service
	auditLogService            AuditLogService
	gatewayService             GatewayService
	guildService               GuildService
	memberService              MemberService
	channelService             ChannelService
	threadService              ThreadService
	interactionService         InteractionService
	inviteService              InviteService
	guildTemplateService       GuildTemplateService
	userService                UserService
	voiceService               VoiceService
	webhookService             WebhookService
	stageInstanceService       StageInstanceService
	emojiService               EmojiService
	stickerService             StickerService
	guildScheduledEventService GuildScheduledEventService
}

func (s *servicesImpl) RestClient() Client {
	return s.restClient
}

func (s *servicesImpl) HTTPClient() *http.Client {
	return s.RestClient().HTTPClient()
}

func (s *servicesImpl) Close(ctx context.Context) {
	s.restClient.Close(ctx)
}

func (s *servicesImpl) Application() ApplicationService {
	return s.applicationService
}

func (s *servicesImpl) OAuth2() OAuth2Service {
	return s.oauth2Service
}

func (s *servicesImpl) AuditLog() AuditLogService {
	return s.auditLogService
}

func (s *servicesImpl) Gateway() GatewayService {
	return s.gatewayService
}

func (s *servicesImpl) Guild() GuildService {
	return s.guildService
}

func (s *servicesImpl) Member() MemberService {
	return s.memberService
}

func (s *servicesImpl) Channel() ChannelService {
	return s.channelService
}

func (s *servicesImpl) Thread() ThreadService {
	return s.threadService
}

func (s *servicesImpl) Interaction() InteractionService {
	return s.interactionService
}

func (s *servicesImpl) Invite() InviteService {
	return s.inviteService
}

func (s *servicesImpl) GuildTemplate() GuildTemplateService {
	return s.guildTemplateService
}

func (s *servicesImpl) User() UserService {
	return s.userService
}

func (s *servicesImpl) Voice() VoiceService {
	return s.voiceService
}

func (s *servicesImpl) Webhook() WebhookService {
	return s.webhookService
}

func (s *servicesImpl) StageInstance() StageInstanceService {
	return s.stageInstanceService
}

func (s *servicesImpl) Emoji() EmojiService {
	return s.emojiService
}

func (s *servicesImpl) Sticker() StickerService {
	return s.stickerService
}

func (s *servicesImpl) GuildScheduledEvent() GuildScheduledEventService {
	return s.guildScheduledEventService
}

type Service interface {
	RestClient() Client
}
