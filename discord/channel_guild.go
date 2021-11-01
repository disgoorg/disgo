package discord

type GuildChannel interface {
	guildChannel()
}

type CategoryChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id"`
	Position               int                   `json:"position"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

type StoreChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id"`
	Position               int                   `json:"position"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}
