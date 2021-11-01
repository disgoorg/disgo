package discord

type GuildVoiceChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id,omitempty"`
	Position               int                   `json:"position,omitempty"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name,omitempty"`
	Topic                  *string               `json:"topic,omitempty"`
	Bitrate                int                   `json:"bitrate,omitempty"`
	UserLimit              int                   `json:"user_limit,omitempty"`
	ParentID               *Snowflake            `json:"parent_id,omitempty"`
	RTCRegion              string                `json:"rtc_region"`
	VideoQualityMode       VideoQualityMode      `json:"video_quality_mode"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

type GuildStageChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id,omitempty"`
	Position               int                   `json:"position,omitempty"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name,omitempty"`
	Topic                  *string               `json:"topic,omitempty"`
	Bitrate                int                   `json:"bitrate,omitempty"`
	UserLimit              int                   `json:"user_limit,omitempty"`
	ParentID               *Snowflake            `json:"parent_id,omitempty"`
	RTCRegion              string                `json:"rtc_region"`
	VideoQualityMode       VideoQualityMode      `json:"video_quality_mode"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

// VideoQualityMode https://discord.com/developers/docs/resources/channel#channel-object-video-quality-modes
type VideoQualityMode int

//goland:noinspection GoUnusedConst
const (
	VideoQualityModeAuto = iota + 1
	VideoQualityModeFull
)
