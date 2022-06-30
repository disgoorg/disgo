package rpc

import (
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

type CmdArgsSetActivity struct {
	PID      int              `json:"pid"`
	Activity discord.Activity `json:"activity"`
}

func (CmdArgsSetActivity) cmdArgs() {}

type CmdRsSetActivity struct {
	discord.Activity
}

func (CmdRsSetActivity) messageData() {}

type CmdArgsSubscribe struct {
	Evt string `json:"evt"`
}

type CmdArgsUnsubscribe struct {
}
