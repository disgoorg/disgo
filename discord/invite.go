package discord

import (
	"bytes"
	"encoding/csv"
	"mime/multipart"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

// InviteTargetType is type of target an Invite uses
type InviteTargetType int

// Constants for TargetType
const (
	InviteTargetTypeStream InviteTargetType = iota + 1
	InviteTargetTypeEmbeddedApplication
	InviteTargetTypeRoleSubscriptionsPurchase
)

// Invite is a partial invite struct
type Invite struct {
	Type                     InviteType           `json:"type"`
	Code                     string               `json:"code"`
	Guild                    *InviteGuild         `json:"guild"`
	Channel                  *InviteChannel       `json:"channel"`
	Inviter                  *User                `json:"inviter"`
	TargetUser               *User                `json:"target_user"`
	TargetType               InviteTargetType     `json:"target_user_type"`
	ApproximatePresenceCount int                  `json:"approximate_presence_count"`
	ApproximateMemberCount   int                  `json:"approximate_member_count"`
	ExpiresAt                *time.Time           `json:"expires_at"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event"`
	Flags                    InviteFlags          `json:"flags"`
	Roles                    []InviteRole         `json:"roles"`
}

func (i Invite) URL() string {
	return InviteURL(i.Code)
}

type InviteType int

const (
	InviteTypeGuild InviteType = iota
	InviteTypeGroupDM
	InviteTypeFriend
)

type InviteFlags int

const (
	InviteFlagIsGuestInvite InviteFlags = 1 << iota
	InviteFlagsNone         InviteFlags = 0
)

type ExtendedInvite struct {
	Invite
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}

type PartialInvite struct {
	Code *string `json:"code"`
	Uses int     `json:"uses"`
}

type InviteChannel struct {
	ID   snowflake.ID `json:"id"`
	Type ChannelType  `json:"type"`
	Name string       `json:"name"`
	Icon *string      `json:"icon,omitempty"`
}

type InviteRole struct {
	ID         snowflake.ID `json:"id"`
	Name       string       `json:"name"`
	Color      int          `json:"color"`
	RoleColors RoleColors   `json:"colors"`
	Position   int          `json:"position"`
	Icon       *string      `json:"icon"`
	Emoji      *string      `json:"unicode_emoji"`
}

// IconURL returns the Icon URL of this channel.
// This will be nil for every ChannelType except ChannelTypeGroupDM
func (c InviteChannel) IconURL(opts ...CDNOpt) *string {
	if c.Icon == nil {
		return nil
	}
	url := formatAssetURL(ChannelIcon, opts, c.ID, *c.Icon)
	return &url
}

// An InviteGuild is the Guild of an Invite
type InviteGuild struct {
	ID                snowflake.ID      `json:"id"`
	Name              string            `json:"name"`
	Splash            *string           `json:"splash"`
	Banner            *string           `json:"banner"`
	Description       *string           `json:"description"`
	Icon              *string           `json:"icon"`
	Features          []GuildFeature    `json:"features"`
	VerificationLevel VerificationLevel `json:"verification_level"`
	VanityURLCode     *string           `json:"vanity_url_code"`
}

func (g InviteGuild) IconURL(opts ...CDNOpt) *string {
	if g.Icon == nil {
		return nil
	}
	url := formatAssetURL(GuildIcon, opts, g.ID, *g.Icon)
	return &url
}

func (g InviteGuild) SplashURL(opts ...CDNOpt) *string {
	if g.Splash == nil {
		return nil
	}
	url := formatAssetURL(GuildSplash, opts, g.ID, *g.Splash)
	return &url
}

type InviteCreate struct {
	MaxAge              *int             `json:"max_age,omitempty"`
	MaxUses             *int             `json:"max_uses,omitempty"`
	Temporary           bool             `json:"temporary,omitempty"`
	Unique              bool             `json:"unique,omitempty"`
	TargetType          InviteTargetType `json:"target_type,omitempty"`
	TargetUserID        snowflake.ID     `json:"target_user_id,omitempty"`
	TargetApplicationID snowflake.ID     `json:"target_application_id,omitempty"`
	TargetUserIDs       []snowflake.ID   `json:"-"`
	RoleIDs             []snowflake.ID   `json:"role_ids,omitempty"`
}

// ToBody returns the InviteCreate ready for body
func (i InviteCreate) ToBody() (any, error) {
	if len(i.TargetUserIDs) > 0 {
		return payloadWithTargetUserIDs(i, i.TargetUserIDs)
	}
	return i, nil
}

type InviteTargetUsersUpdate []snowflake.ID

// ToBody returns the InviteTargetUsersUpdate ready for body
func (i InviteTargetUsersUpdate) ToBody() (any, error) {
	return payloadWithTargetUserIDs(nil, i)
}

func payloadWithTargetUserIDs(v any, targetUsersIDs []snowflake.ID) (*MultipartBuffer, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	defer func() {
		_ = writer.Close()
	}()

	if v != nil {
		payload, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreatePart(partHeader(`form-data; name="payload_json"`, "application/json"))
		if err != nil {
			return nil, err
		}

		if _, err = part.Write(payload); err != nil {
			return nil, err
		}
	}

	part, err := writer.CreatePart(partHeader(`form-data; name="target_users_file"; filename="users.csv"`, "text/csv"))
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(part)
	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"Users"}); err != nil {
		return nil, err
	}

	for _, userID := range targetUsersIDs {
		if err := csvWriter.Write([]string{userID.String()}); err != nil {
			return nil, err
		}
	}

	return &MultipartBuffer{
		Buffer:      buffer,
		ContentType: writer.FormDataContentType(),
	}, nil
}

type TargetUsersJobStatus struct {
	Status         TargetUsersJobStatusCode `json:"status"`
	TotalUsers     int                      `json:"total_users"`
	ProcessedUsers int                      `json:"processed_users"`
	CreatedAt      time.Time                `json:"created_at"`
	CompletedAt    *time.Time               `json:"completed_at"`
	ErrorMessage   *string                  `json:"error_message"`
}

// TargetUsersJobStatusCode indicates the status of creating or updating target users of invite
type TargetUsersJobStatusCode int

// all TargetUsersJobStatusCode
const (
	TargetUsersJobStatusCodeUnspecified TargetUsersJobStatusCode = iota
	TargetUsersJobStatusCodeProcessing
	TargetUsersJobStatusCodeCompleted
	TargetUsersJobStatusCodeFailed
)
