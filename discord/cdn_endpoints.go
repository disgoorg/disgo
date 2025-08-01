package discord

import (
	"strings"
)

const (
	CDN      = "https://cdn.discordapp.com"
	CDNMedia = "https://media.discordapp.net"
)

var (
	CustomEmoji = NewCDN("/emojis/{emote.id}", FileFormatPNG, FileFormatWebP, FileFormatAVIF, FileFormatGIF)

	GuildIcon            = NewCDN("/icons/{guild.id}/{guild.icon.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP, FileFormatGIF)
	GuildSplash          = NewCDN("/splashes/{guild.id}/{guild.splash.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)
	GuildDiscoverySplash = NewCDN("/discovery-splashes/{guild.id}/{guild.discovery.splash.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)
	GuildBanner          = NewCDN("/banners/{guild.id}/{guild.banner.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP, FileFormatGIF)

	GuildScheduledEventCover = NewCDN("/guild-events/{event.id}/{event.cover.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	RoleIcon = NewCDN("/role-icons/{role.id}/{role.icon.hash}", FileFormatPNG, FileFormatJPEG)

	UserBanner        = NewCDN("/banners/{user.id}/{user.banner.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP, FileFormatGIF)
	UserAvatar        = NewCDN("/avatars/{user.id}/{user.avatar.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP, FileFormatGIF)
	DefaultUserAvatar = NewCDN("/embed/avatars/{index}", FileFormatPNG)

	ChannelIcon = NewCDN("/channel-icons/{channel.id}/{channel.icon.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	MemberAvatar = NewCDN("/guilds/{guild.id}/users/{user.id}/avatars/{member.avatar.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP, FileFormatGIF)
	MemberBanner = NewCDN("/guilds/{guild.id}/users/{user.id}/banners/{member.avatar.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP, FileFormatGIF)

	AvatarDecoration = NewCDN("/avatar-decoration-presets/{user.avatar.decoration.hash}", FileFormatPNG)

	GuildTagBadge = NewCDN("/guild-tag-badges/{guild.id}/{badge.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	ApplicationIcon  = NewCDN("/app-icons/{application.id}/{icon.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)
	ApplicationCover = NewCDN("/app-assets/{application.id}/{cover.image.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)
	ApplicationAsset = NewCDN("/app-assets/{application.id}/{asset.id}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	AchievementIcon = NewCDN("/app-assets/{application.id}/achievements/{achievement.id}/icons/{icon.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	StorePageAsset = NewCDN("/app-assets/{application.id}/store/{asset.id}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	TeamIcon = NewCDN("/team-icons/{team.id}/{team.icon.hash}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)

	StickerPackBanner = NewCDN("/app-assets/710982414301790216/store/{banner.asset.id}", FileFormatPNG, FileFormatJPEG, FileFormatWebP)
	CustomSticker     = NewCDN("/stickers/{sticker.id}", FileFormatPNG, FileFormatLottie, FileFormatGIF)

	AttachmentFile = NewCDN("/attachments/{channel.id}/{attachment.id}/{file.name}", FileFormatNone)

	SoundboardSoundFile = NewCDN("/soundboard-sounds/{sound.id}", FileFormatNone)

	NameplateAsset = NewCDN("/assets/collectibles/{nameplate.asset}/asset", FileFormatWebM)
)

// FileFormat is the type of file on Discord's CDN (https://discord.com/developers/docs/reference#image-formatting-image-formats)
type FileFormat string

// The available FileFormat(s)
const (
	FileFormatNone   FileFormat = ""
	FileFormatPNG    FileFormat = "png"
	FileFormatJPEG   FileFormat = "jpg"
	FileFormatWebP   FileFormat = "webp"
	FileFormatAVIF   FileFormat = "avif"
	FileFormatGIF    FileFormat = "gif"
	FileFormatLottie FileFormat = "json"
	FileFormatWebM   FileFormat = "webm"
)

// String returns the string representation of the FileFormat
func (f FileFormat) String() string {
	return string(f)
}

// Animated returns true if the FileFormat is animated
func (f FileFormat) Animated() bool {
	switch f {
	case FileFormatWebP, FileFormatAVIF, FileFormatGIF, FileFormatWebM:
		return true
	default:
		return false
	}
}

func NewCDN(route string, fileFormats ...FileFormat) *CDNEndpoint {
	return &CDNEndpoint{
		Route:   route,
		Formats: fileFormats,
	}
}

type CDNEndpoint struct {
	Route   string
	Formats []FileFormat
}

func (e CDNEndpoint) URL(format FileFormat, values QueryValues, params ...any) string {
	query := values.Encode()
	if query != "" {
		query = "?" + query
	}

	// for some reason custom gif stickers use a different cdn url, blame discord for this one
	if format == FileFormatGIF && e.Route == "/stickers/{sticker.id}" {
		return urlPrint(CDNMedia+e.Route+"."+format.String(), params...) + query
	}
	route := CDN + e.Route
	// only append period and file extension if the format is not FileFormatNone
	if format != FileFormatNone {
		route += "." + format.String()
	}

	return urlPrint(route, params...) + query
}

func DefaultCDNConfig(format FileFormat) *CDNConfig {
	return &CDNConfig{
		Format: format,
		Values: QueryValues{},
	}
}

type CDNConfig struct {
	Format FileFormat
	Values QueryValues
}

// Apply applies the given ConfigOpt(s) to the Config
func (c *CDNConfig) Apply(opts []CDNOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

type CDNOpt func(config *CDNConfig)

func WithSize(size int) CDNOpt {
	return func(config *CDNConfig) {
		config.Values["size"] = size
	}
}

func WithFormat(format FileFormat) CDNOpt {
	return func(config *CDNConfig) {
		config.Format = format
	}
}

func formatAssetURL(cdnRoute *CDNEndpoint, opts []CDNOpt, params ...any) string {
	format := FileFormatNone
	if len(cdnRoute.Formats) > 0 { // just in case someone fucks up
		// use the first provided format in the route definition itself. if the user provides a different format, this will be overriden by the apply function call below
		// previously, the default format was png, which would cause issues for cdn endpoints like attachments and soundboard sounds, requiring custom "overrides"
		format = cdnRoute.Formats[0]
	}
	cfg := DefaultCDNConfig(format)
	cfg.Apply(opts)

	var lastStringParam string
	lastParam := params[len(params)-1]
	if str, ok := lastParam.(string); ok {
		if str == "" {
			return ""
		}
		lastStringParam = str
	} else if ptrStr, ok := lastParam.(*string); ok {
		if ptrStr == nil {
			return ""
		}
		lastStringParam = *ptrStr
	}

	// some endpoints have a_ prefix for animated images except the AvatarDecoration endpoint does not like this
	if strings.HasPrefix(lastStringParam, "a_") && !cfg.Format.Animated() && cdnRoute.Route != "/avatar-decoration-presets/{user.avatar.decoration.hash}" {
		cfg.Format = FileFormatGIF
	}

	return cdnRoute.URL(cfg.Format, cfg.Values, params...)
}
