package discord

import (
	"errors"
)

var (
	ErrNoGatewayOrShardManager = errors.New("no gateway or shard manager configured")
	ErrNoGuildMembersIntent    = errors.New("this operation requires the GUILD_MEMBERS intent")
	ErrNoShardManager          = errors.New("no shard manager configured")
	ErrNoGateway               = errors.New("no gateway configured")
	ErrGatewayAlreadyConnected = errors.New("gateway is already connected")
	ErrShardNotConnected       = errors.New("shard is not connected")
	ErrShardNotFound           = errors.New("shard not found in shard manager")
	ErrNoHTTPServer            = errors.New("no http server configured")

	ErrInvalidBotToken = errors.New("token is not in a valid format")
	ErrNoBotToken      = errors.New("please specify the token")

	ErrInteractionAlreadyReplied = errors.New("you already replied to this interaction")
	ErrInteractionExpired        = errors.New("this interaction has expired")

	ErrCheckFailed = errors.New("check failed")
)
