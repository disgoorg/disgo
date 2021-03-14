package api

// Role is a Guild Role object
type Role struct {
	ID Snowflake
	Name string
	GuildID Snowflake
}
// todo: add other props, missing several