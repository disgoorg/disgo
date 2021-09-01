package discord

import (
	"errors"
	"fmt"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrBadGateway   = errors.New("bad gateway could not reach destination")
	ErrUnauthorized = errors.New("not authorized for this endpoint")
	ErrBadRequest   = errors.New("bad request")
	ErrRatelimited  = errors.New("received error 429")

	ErrNoGateway             = errors.New("no gateway configured")
	ErrNoGatewayConn         = errors.New("gateway is not connected")
	ErrGatewayCompressedData = errors.New("disgo does not currently support compressed gateway data")

	ErrNoDisgoInstance = errors.New("no disgo instance injected")

	ErrInvalidBotToken = errors.New("BotToken is not in a valid format")
	ErrNoBotToken      = errors.New("please specify the BotToken")

	ErrSelfDM = errors.New("can't open a dm channel to yourself")

	ErrInteractionAlreadyReplied = errors.New("you already replied to this interaction")

	ErrChannelNotTypeNews = errors.New("channel type is not 'NEWS'")

	ErrCheckFailed = errors.New("check failed")
)

func ErrUnexpectedQueryParam(param string) error {
	return fmt.Errorf("unexpected query param '%s' received", param)
}

func ErrInvalidArgCount(argCount int, paramCount int) error {
	return fmt.Errorf("invalid amount of arguments received. expected: %d, received: %d", argCount, paramCount)
}

func ErrFileExtensionNotSupported(fileExtension string) error {
	return fmt.Errorf("provided file extension: %s is not supported by discord on this end", fileExtension)
}

func ErrUnexpectedGatewayOp(wOp Op, rOp int) error {
	return fmt.Errorf("expected op: %d, received: %d", wOp, rOp)
}
