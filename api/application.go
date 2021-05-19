package api

type ApplicationFlags int

const (
	ApplicationFlagGatewayPresence ApplicationFlags = 1 << (iota + 12)
	ApplicationFlagGatewayPresenceLimited
	applicationFlagGatewayGuildMembers
	ApplicationFlagGatewayGuildMembersLimited
	ApplicationFlagVerificationPendingGuildLimit
	applicationFlagEmbedded
)

type Application struct {
	ID                  Snowflake        `json:"id"`
	Name                string           `json:"name"`
	Icon                *string          `json:"icon"`
	Description         string           `json:"description"`
	RPCOrigins          []string         `json:"rpc_origins"`
	BotPublic           bool             `json:"bot_public"`
	BotRequireCodeGrant bool             `json:"bot_require_code_grant"`
	TermsOfServiceURL   *string          `json:"terms_of_service_url"`
	PrivacyPolicyURL    *string          `json:"privacy_policy_url"`
	Owner               *User            `json:"owner"`
	Summary             string           `json:"summary"`
	VerifyKey           string           `json:"verify_key"`
	Team                *ApplicationTeam `json:"team"`
	GuildID             *Snowflake       `json:"guild_id"`
	PrimarySkuID        Snowflake        `json:"primary_sku_id"`
	Slug                *string          `json:"slug"`
	CoverImage          string           `json:"cover_image"`
	Flags               ApplicationFlags `json:"flags"`
}

type ApplicationTeam struct {
	Icon        string               `json:"icon"`
	ID          Snowflake            `json:"id"`
	Members     []*ApplicationMember `json:"members"`
	Name        string               `json:"name"`
	OwnerUserID Snowflake            `json:"owner_user_id"`
}

type MembershipState int

const (
	MembershipStateInvited MembershipState = iota + 1
	MembershipStateAccepted
)

type ApplicationMember struct {
	MembershipState MembershipState `json:"membership_state"`
	Permissions     []string        `json:"permissions"`
	TeamID          Snowflake       `json:"team_id"`
	User            *User           `json:"user"`
}
