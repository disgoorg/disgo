package discord

import (
	"github.com/pkg/errors"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrNoGatewayOrShardManager = errors.New("no gateway or shard manager configured")
	ErrNoGuildMembersIntent    = errors.New("this operation requires the GUILD_MEMBERS intent")
	ErrNoShardManager          = errors.New("no shard manager configured")
	ErrNoGateway               = errors.New("no gateway configured")
	ErrShardNotConnected       = errors.New("shard is not connected")
	ErrShardNotFound           = errors.New("shard not found in shard manager")
	ErrGatewayCompressedData   = errors.New("disgo does not currently support compressed gateway data")
	ErrNoHTTPServer            = errors.New("no http server configured")

	ErrNoDisgoInstance = errors.New("no disgo instance injected")

	ErrInvalidBotToken = errors.New("BotToken is not in a valid format")
	ErrNoBotToken      = errors.New("please specify the BotToken")

	ErrSelfDM = errors.New("can't open a dm channel to yourself")

	ErrInteractionAlreadyReplied = errors.New("you already replied to this interaction")

	ErrChannelNotTypeNews = errors.New("channel type is not 'NEWS'")

	ErrCheckFailed = errors.New("check failed")

	ErrMemberMustBeConnectedToChannel = errors.New("the member must be connected to the channel")

	ErrStickerTypeGuild = errors.New("sticker type must be of type StickerTypeGuild")
)
