package api

type Emote struct {
	ID   Snowflake
	Name string
}

func (e Emote) Mention() string {
	return "<:" + e.Name +  ":" + e.ID.String() + ">"
}

func (e Emote) Reaction() string {
	return ":" + e.Name +  ":" + e.ID.String()
}

