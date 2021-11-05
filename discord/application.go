package discord

import (
	"fmt"
	"strings"
)

type Application struct {
	ID                  Snowflake        `json:"id"`
	Name                string           `json:"name"`
	Icon                *string          `json:"icon,omitempty"`
	Description         string           `json:"description"`
	RPCOrigins          []string         `json:"rpc_origins"`
	BotPublic           bool             `json:"bot_public"`
	BotRequireCodeGrant bool             `json:"bot_require_code_grant"`
	TermsOfServiceURL   *string          `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL    *string          `json:"privacy_policy_url,omitempty"`
	Owner               *User            `json:"owner,omitempty"`
	Summary             string           `json:"summary"`
	VerifyKey           string           `json:"verify_key"`
	Team                *Team            `json:"team,omitempty"`
	GuildID             *Snowflake       `json:"guild_id,omitempty"`
	PrimarySkuID        *Snowflake       `json:"primary_sku_id,omitempty"`
	Slug                *string          `json:"slug,omitempty"`
	CoverImage          *string          `json:"cover_image,omitempty"`
	Flags               ApplicationFlags `json:"flags,omitempty"`
}

type PartialApplication struct {
	ID    Snowflake        `json:"id"`
	Flags ApplicationFlags `json:"flags"`
}

type AuthorizationInformation struct {
	Application Application        `json:"application"`
	Scopes      []ApplicationScope `json:"scopes"`
	Expires     Time               `json:"expires"`
	User        *User              `json:"user"`
}

type ApplicationScope string

//goland:noinspection GoUnusedConst
const (
	ApplicationScopeActivitiesWrite ApplicationScope = "activities.write"
	ApplicationScopeActivitiesRead  ApplicationScope = "activities.read"

	ApplicationScopeRPC                  ApplicationScope = "rpc"
	ApplicationScopeRPCNotificationsRead ApplicationScope = "rpc.notifications.read"
	ApplicationScopeRPCVoiceWrite        ApplicationScope = "rpc.voice.write"
	ApplicationScopeRPCVoiceRead         ApplicationScope = "rpc.voice.read"
	ApplicationScopeRPCActivitiesWrite   ApplicationScope = "rpc.activities.write"

	ApplicationScopeGuilds     ApplicationScope = "guilds"
	ApplicationScopeGuildsJoin ApplicationScope = "guilds.join"
	ApplicationScopeGDMJoin    ApplicationScope = "gdm.join"

	ApplicationScopeRelationshipsRead ApplicationScope = "relationships.read"
	ApplicationScopeIdentify          ApplicationScope = "identify"
	ApplicationScopeEmail             ApplicationScope = "email"
	ApplicationScopeConnections       ApplicationScope = "connections"
	ApplicationScopeBot               ApplicationScope = "bot"
	ApplicationScopeMessagesRead      ApplicationScope = "messages.read"
	ApplicationScopeWebhookIncoming   ApplicationScope = "webhook.incoming"

	ApplicationScopeApplicationsCommands       ApplicationScope = "applications.commands"
	ApplicationScopeApplicationsCommandsUpdate ApplicationScope = "applications.commands.update"
	ApplicationScopeApplicationsEntitlements   ApplicationScope = "applications.entitlements"
	ApplicationScopeApplicationsStoreUpdate    ApplicationScope = "applications.store.update"
	ApplicationScopeApplicationsBuildsRead     ApplicationScope = "applications.builds.read"
	ApplicationScopeApplicationsBuildsUpload   ApplicationScope = "applications.builds.upload"
)

func (s ApplicationScope) String() string {
	return string(s)
}

const ScopeSeparator = " "

func JoinScopes(scopes []ApplicationScope) string {
	strScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		strScopes[i] = scope.String()
	}
	return strings.Join(strScopes, ScopeSeparator)
}

func SplitScopes(joinedScopes string) []ApplicationScope {
	var scopes []ApplicationScope
	for _, scope := range strings.Split(joinedScopes, ScopeSeparator) {
		scopes = append(scopes, ApplicationScope(scope))
	}
	return scopes
}

func HasScope(scope ApplicationScope, scopes ...ApplicationScope) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

type TokenType string

//goland:noinspection GoUnusedConst
const (
	TokenTypeBearer TokenType = "Bearer"
	TokenTypeBot    TokenType = "Bot"
)

func (t TokenType) String() string {
	return string(t)
}

func (t TokenType) Apply(token string) string {
	return fmt.Sprintf("%s %s", t.String(), token)
}

// ApplicationFlags (https://discord.com/developers/docs/resources/application#application-object-application-flags)
type ApplicationFlags int

//goland:noinspection GoUnusedConst
const (
	ApplicationFlagGatewayPresence = 1 << (iota + 12)
	ApplicationFlagGatewayPresenceLimited
	ApplicationFlagGatewayGuildMembers
	ApplicationFlagGatewayGuildMemberLimited
	ApplicationFlagVerificationPendingGuildLimit
	ApplicationFlagEmbedded
)

// Add allows you to add multiple bits together, producing a new bit
func (f ApplicationFlags) Add(bits ...ApplicationFlags) ApplicationFlags {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f ApplicationFlags) Remove(bits ...ApplicationFlags) ApplicationFlags {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func (f ApplicationFlags) Has(bits ...ApplicationFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (f ApplicationFlags) Missing(bits ...ApplicationFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}

type Team struct {
	Icon    *string      `json:"icon"`
	ID      string       `json:"id"`
	Members []TeamMember `json:"members"`
	Name    string       `json:"name"`
	OwnerID Snowflake    `json:"owner_user_id"`
}

type TeamMember struct {
	MembershipState MembershipState   `json:"membership_state"`
	Permissions     []TeamPermissions `json:"permissions"`
	TeamID          Snowflake         `json:"team_id"`
	User            User              `json:"user"`
}

type MembershipState int

//goland:noinspection GoUnusedConst,GoUnusedConst
const (
	MembershipStateInvited = iota + 1
	MembershipStateAccepted
)

type TeamPermissions string

//goland:noinspection GoUnusedConst
const (
	TeamPermissionAdmin = "*"
)
