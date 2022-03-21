package discord

import (
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

var _ Interaction = (*PingInteraction)(nil)

type PingInteraction struct {
	id            snowflake.Snowflake
	applicationID snowflake.Snowflake
	token         string
	version       int
}

func (i *PingInteraction) UnmarshalJSON(data []byte) error {
	var v rawInteraction

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	i.id = v.ID
	i.applicationID = v.ApplicationID
	i.token = v.Token
	i.version = v.Version
	return nil
}

func (i PingInteraction) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawInteraction{
		ID:            i.id,
		Type:          i.Type(),
		ApplicationID: i.applicationID,
		Token:         i.token,
		Version:       i.version,
	})
}

func (PingInteraction) Type() InteractionType {
	return InteractionTypePing
}

func (i PingInteraction) ID() snowflake.Snowflake {
	return i.id
}

func (i PingInteraction) ApplicationID() snowflake.Snowflake {
	return i.applicationID
}

func (i PingInteraction) Token() string {
	return i.token
}

func (i PingInteraction) Version() int {
	return i.version
}

func (PingInteraction) GuildID() *snowflake.Snowflake {
	return nil
}

func (PingInteraction) ChannelID() snowflake.Snowflake {
	return ""
}

func (PingInteraction) Locale() Locale {
	return ""
}

func (PingInteraction) GuildLocale() *Locale {
	return nil
}

func (PingInteraction) Member() *ResolvedMember {
	return nil
}

func (PingInteraction) User() User {
	return User{}
}

func (PingInteraction) interaction() {}
