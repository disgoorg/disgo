package rpc

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type CmdArgsAuthorize struct {
	ClientID snowflake.ID          `json:"client_id"`
	Scopes   []discord.OAuth2Scope `json:"scopes"`
	RPCToken string                `json:"rpc_token,omitempty"`
	Username string                `json:"username,omitempty"`
}

func (CmdArgsAuthorize) cmdArgs() {}

type CmdRsAuthorize struct {
	Code string `json:"code"`
}

func (CmdRsAuthorize) messageData() {}

type CmdArgsAuthenticate struct {
	AccessToken string `json:"access_token"`
}

func (CmdArgsAuthenticate) cmdArgs() {}

type CmdRsAuthenticate struct {
	User        discord.User              `json:"user"`
	Scopes      []discord.OAuth2Scope     `json:"scopes"`
	Expires     time.Time                 `json:"expires"`
	Application discord.OAuth2Application `json:"application"`
}

func (CmdRsAuthenticate) messageData() {}

type CmdArgsSetActivity struct {
	PID      int              `json:"pid"`
	Activity discord.Activity `json:"activity"`
}

func (CmdArgsSetActivity) cmdArgs() {}

type CmdRsSetActivity struct {
	discord.Activity
}

func (CmdRsSetActivity) messageData() {}

type CmdArgsSubscribe interface {
	CmdArgs
	cmdArgsSubscribe()
}

type CmdArgsSubscribeMessage struct {
	ChannelID snowflake.ID `json:"channel_id"`
}

func (CmdArgsSubscribeMessage) cmdArgs()          {}
func (CmdArgsSubscribeMessage) cmdArgsSubscribe() {}

type CmdRsSubscribe struct {
	Evt string `json:"evt"`
}

func (CmdRsSubscribe) messageData() {}

type CmdArgsUnsubscribe struct {
}

func (CmdArgsUnsubscribe) cmdArgs() {}

type CmdRsUnsubscribe struct {
	Evt string `json:"evt"`
}

func (CmdRsUnsubscribe) messageData() {}
