package rest

import (
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/disgoorg/json/v2"
)

// JSONErrorCode is the error code returned by the Discord API.
// See https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
type JSONErrorCode int

const (
	// General error (such as a malformed request body, amongst other things)
	JSONErrorCodeGeneral JSONErrorCode = 0

	// Unknown resources
	JSONErrorCodeUnknownAccount                        JSONErrorCode = 10001
	JSONErrorCodeUnknownApplication                    JSONErrorCode = 10002
	JSONErrorCodeUnknownChannel                        JSONErrorCode = 10003
	JSONErrorCodeUnknownGuild                          JSONErrorCode = 10004
	JSONErrorCodeUnknownIntegration                    JSONErrorCode = 10005
	JSONErrorCodeUnknownInvite                         JSONErrorCode = 10006
	JSONErrorCodeUnknownMember                         JSONErrorCode = 10007
	JSONErrorCodeUnknownMessage                        JSONErrorCode = 10008
	JSONErrorCodeUnknownPermissionOverwrite            JSONErrorCode = 10009
	JSONErrorCodeUnknownProvider                       JSONErrorCode = 10010
	JSONErrorCodeUnknownRole                           JSONErrorCode = 10011
	JSONErrorCodeUnknownToken                          JSONErrorCode = 10012
	JSONErrorCodeUnknownUser                           JSONErrorCode = 10013
	JSONErrorCodeUnknownEmoji                          JSONErrorCode = 10014
	JSONErrorCodeUnknownWebhook                        JSONErrorCode = 10015
	JSONErrorCodeUnknownWebhookService                 JSONErrorCode = 10016
	JSONErrorCodeUnknownSession                        JSONErrorCode = 10020
	JSONErrorCodeUnknownAsset                          JSONErrorCode = 10021
	JSONErrorCodeUnknownBan                            JSONErrorCode = 10026
	JSONErrorCodeUnknownSKU                            JSONErrorCode = 10027
	JSONErrorCodeUnknownStoreListing                   JSONErrorCode = 10028
	JSONErrorCodeUnknownEntitlement                    JSONErrorCode = 10029
	JSONErrorCodeUnknownBuild                          JSONErrorCode = 10030
	JSONErrorCodeUnknownLobby                          JSONErrorCode = 10031
	JSONErrorCodeUnknownBranch                         JSONErrorCode = 10032
	JSONErrorCodeUnknownStoreDirectoryLayout           JSONErrorCode = 10033
	JSONErrorCodeUnknownRedistributable                JSONErrorCode = 10036
	JSONErrorCodeUnknownGiftCode                       JSONErrorCode = 10038
	JSONErrorCodeUnknownStream                         JSONErrorCode = 10049
	JSONErrorCodeUnknownPremiumServerSubscribeCooldown JSONErrorCode = 10050
	JSONErrorCodeUnknownGuildTemplate                  JSONErrorCode = 10057
	JSONErrorCodeUnknownDiscoverableServerCategory     JSONErrorCode = 10059
	JSONErrorCodeUnknownSticker                        JSONErrorCode = 10060
	JSONErrorCodeUnknownStickerPack                    JSONErrorCode = 10061
	JSONErrorCodeUnknownInteraction                    JSONErrorCode = 10062
	JSONErrorCodeUnknownApplicationCommand             JSONErrorCode = 10063
	JSONErrorCodeUnknownVoiceState                     JSONErrorCode = 10065
	JSONErrorCodeUnknownApplicationCommandPermissions  JSONErrorCode = 10066
	JSONErrorCodeUnknownStageInstance                  JSONErrorCode = 10067
	JSONErrorCodeUnknownGuildMemberVerificationForm    JSONErrorCode = 10068
	JSONErrorCodeUnknownGuildWelcomeScreen             JSONErrorCode = 10069
	JSONErrorCodeUnknownGuildScheduledEvent            JSONErrorCode = 10070
	JSONErrorCodeUnknownGuildScheduledEventUser        JSONErrorCode = 10071
	JSONErrorCodeUnknownTag                            JSONErrorCode = 10087
	JSONErrorCodeUnknownSound                          JSONErrorCode = 10097

	// Authorization/action errors
	JSONErrorCodeBotsCannotUseThisEndpoint                   JSONErrorCode = 20001
	JSONErrorCodeOnlyBotsCanUseThisEndpoint                  JSONErrorCode = 20002
	JSONErrorCodeExplicitContentCannotBeSent                 JSONErrorCode = 20009
	JSONErrorCodeNotAuthorizedToPerformThisAction            JSONErrorCode = 20012
	JSONErrorCodeActionCannotBePerformedDueToSlowmode        JSONErrorCode = 20016
	JSONErrorCodeOnlyOwnerCanPerformThisAction               JSONErrorCode = 20018
	JSONErrorCodeMessageCannotBeEditedDueToAnnouncementLimit JSONErrorCode = 20022
	JSONErrorCodeUnderMinimumAge                             JSONErrorCode = 20024
	JSONErrorCodeChannelWriteRateLimit                       JSONErrorCode = 20028
	JSONErrorCodeServerWriteRateLimit                        JSONErrorCode = 20029
	JSONErrorCodeStageTopicServerNameOrChannelNamesBlocked   JSONErrorCode = 20031
	JSONErrorCodeGuildPremiumSubscriptionLevelTooLow         JSONErrorCode = 20035

	// Maximum limits
	JSONErrorCodeMaximumGuildsReached                          JSONErrorCode = 30001
	JSONErrorCodeMaximumFriendsReached                         JSONErrorCode = 30002
	JSONErrorCodeMaximumPinsReached                            JSONErrorCode = 30003
	JSONErrorCodeMaximumRecipientsReached                      JSONErrorCode = 30004
	JSONErrorCodeMaximumGuildRolesReached                      JSONErrorCode = 30005
	JSONErrorCodeMaximumWebhooksReached                        JSONErrorCode = 30007
	JSONErrorCodeMaximumEmojisReached                          JSONErrorCode = 30008
	JSONErrorCodeMaximumReactionsReached                       JSONErrorCode = 30010
	JSONErrorCodeMaximumGroupDMsReached                        JSONErrorCode = 30011
	JSONErrorCodeMaximumGuildChannelsReached                   JSONErrorCode = 30013
	JSONErrorCodeMaximumAttachmentsInMessageReached            JSONErrorCode = 30015
	JSONErrorCodeMaximumInvitesReached                         JSONErrorCode = 30016
	JSONErrorCodeMaximumAnimatedEmojisReached                  JSONErrorCode = 30018
	JSONErrorCodeMaximumServerMembersReached                   JSONErrorCode = 30019
	JSONErrorCodeMaximumServerCategoriesReached                JSONErrorCode = 30030
	JSONErrorCodeGuildAlreadyHasTemplate                       JSONErrorCode = 30031
	JSONErrorCodeMaximumApplicationCommandsReached             JSONErrorCode = 30032
	JSONErrorCodeMaximumThreadParticipantsReached              JSONErrorCode = 30033
	JSONErrorCodeMaximumDailyApplicationCommandCreatesReached  JSONErrorCode = 30034
	JSONErrorCodeMaximumBansForNonGuildMembersExceeded         JSONErrorCode = 30035
	JSONErrorCodeMaximumBansFetchesReached                     JSONErrorCode = 30037
	JSONErrorCodeMaximumUncompletedGuildScheduledEventsReached JSONErrorCode = 30038
	JSONErrorCodeMaximumStickersReached                        JSONErrorCode = 30039
	JSONErrorCodeMaximumPruneRequestsReached                   JSONErrorCode = 30040
	JSONErrorCodeMaximumGuildWidgetSettingsUpdatesReached      JSONErrorCode = 30042
	JSONErrorCodeMaximumSoundboardSoundsReached                JSONErrorCode = 30045
	JSONErrorCodeMaximumEditsToMessagesOlderThan1HourReached   JSONErrorCode = 30046
	JSONErrorCodeMaximumPinnedThreadsInForumChannelReached     JSONErrorCode = 30047
	JSONErrorCodeMaximumTagsInForumChannelReached              JSONErrorCode = 30048
	JSONErrorCodeBitrateTooHighForChannelType                  JSONErrorCode = 30052
	JSONErrorCodeMaximumPremiumEmojisReached                   JSONErrorCode = 30056
	JSONErrorCodeMaximumWebhooksPerGuildReached                JSONErrorCode = 30058
	JSONErrorCodeMaximumChannelPermissionOverwritesReached     JSONErrorCode = 30060
	JSONErrorCodeChannelsForGuildTooLarge                      JSONErrorCode = 30061

	// Temporary/rate limit errors
	JSONErrorCodeUnauthorized                            JSONErrorCode = 40001
	JSONErrorCodeNeedToVerifyAccount                     JSONErrorCode = 40002
	JSONErrorCodeOpeningDirectMessagesTooFast            JSONErrorCode = 40003
	JSONErrorCodeSendMessagesTemporarilyDisabled         JSONErrorCode = 40004
	JSONErrorCodeRequestEntityTooLarge                   JSONErrorCode = 40005
	JSONErrorCodeFeatureTemporarilyDisabled              JSONErrorCode = 40006
	JSONErrorCodeUserBannedFromGuild                     JSONErrorCode = 40007
	JSONErrorCodeConnectionRevoked                       JSONErrorCode = 40012
	JSONErrorCodeOnlyConsumableSKUsCanBeConsumed         JSONErrorCode = 40018
	JSONErrorCodeCanOnlyDeleteSandboxEntitlements        JSONErrorCode = 40019
	JSONErrorCodeTargetUserNotConnectedToVoice           JSONErrorCode = 40032
	JSONErrorCodeMessageAlreadyCrossposted               JSONErrorCode = 40033
	JSONErrorCodeApplicationCommandWithNameAlreadyExists JSONErrorCode = 40041
	JSONErrorCodeApplicationInteractionFailedToSend      JSONErrorCode = 40043
	JSONErrorCodeCannotSendMessageInForumChannel         JSONErrorCode = 40058
	JSONErrorCodeInteractionAlreadyAcknowledged          JSONErrorCode = 40060
	JSONErrorCodeTagNamesMustBeUnique                    JSONErrorCode = 40061
	JSONErrorCodeServiceResourceRateLimited              JSONErrorCode = 40062
	JSONErrorCodeNoTagsAvailableForNonModerators         JSONErrorCode = 40066
	JSONErrorCodeTagRequiredToCreateForumPost            JSONErrorCode = 40067
	JSONErrorCodeEntitlementAlreadyGrantedForResource    JSONErrorCode = 40074
	JSONErrorCodeInteractionMaxFollowUpMessagesReached   JSONErrorCode = 40094

	// Cloudflare blocking
	JSONErrorCodeCloudflareBlockingRequest JSONErrorCode = 40333

	// Invalid/missing errors
	JSONErrorCodeMissingAccess                                JSONErrorCode = 50001
	JSONErrorCodeInvalidAccountType                           JSONErrorCode = 50002
	JSONErrorCodeCannotExecuteActionOnDMChannel               JSONErrorCode = 50003
	JSONErrorCodeGuildWidgetDisabled                          JSONErrorCode = 50004
	JSONErrorCodeCannotEditMessageAuthoredByAnotherUser       JSONErrorCode = 50005
	JSONErrorCodeCannotSendEmptyMessage                       JSONErrorCode = 50006
	JSONErrorCodeCannotSendMessagesToThisUser                 JSONErrorCode = 50007
	JSONErrorCodeCannotSendMessagesInNonTextChannel           JSONErrorCode = 50008
	JSONErrorCodeChannelVerificationLevelTooHigh              JSONErrorCode = 50009
	JSONErrorCodeOAuth2ApplicationDoesNotHaveBot              JSONErrorCode = 50010
	JSONErrorCodeOAuth2ApplicationLimitReached                JSONErrorCode = 50011
	JSONErrorCodeInvalidOAuth2State                           JSONErrorCode = 50012
	JSONErrorCodeLackPermissionsToPerformAction               JSONErrorCode = 50013
	JSONErrorCodeInvalidAuthenticationToken                   JSONErrorCode = 50014
	JSONErrorCodeNoteTooLong                                  JSONErrorCode = 50015
	JSONErrorCodeTooFewOrTooManyMessagesToDelete              JSONErrorCode = 50016
	JSONErrorCodeInvalidMFALevel                              JSONErrorCode = 50017
	JSONErrorCodeMessageCanOnlyBePinnedToOriginalChannel      JSONErrorCode = 50019
	JSONErrorCodeInviteCodeInvalidOrTaken                     JSONErrorCode = 50020
	JSONErrorCodeCannotExecuteActionOnSystemMessage           JSONErrorCode = 50021
	JSONErrorCodeCannotExecuteActionOnThisChannelType         JSONErrorCode = 50024
	JSONErrorCodeInvalidOAuth2AccessToken                     JSONErrorCode = 50025
	JSONErrorCodeMissingRequiredOAuth2Scope                   JSONErrorCode = 50026
	JSONErrorCodeInvalidWebhookToken                          JSONErrorCode = 50027
	JSONErrorCodeInvalidRole                                  JSONErrorCode = 50028
	JSONErrorCodeInvalidRecipients                            JSONErrorCode = 50033
	JSONErrorCodeMessageTooOldToBulkDelete                    JSONErrorCode = 50034
	JSONErrorCodeInvalidFormBody                              JSONErrorCode = 50035
	JSONErrorCodeInviteAcceptedToGuildBotNotIn                JSONErrorCode = 50036
	JSONErrorCodeInvalidActivityAction                        JSONErrorCode = 50039
	JSONErrorCodeInvalidAPIVersion                            JSONErrorCode = 50041
	JSONErrorCodeFileUploadedExceedsMaximumSize               JSONErrorCode = 50045
	JSONErrorCodeInvalidFileUploaded                          JSONErrorCode = 50046
	JSONErrorCodeCannotSelfRedeemGift                         JSONErrorCode = 50054
	JSONErrorCodeInvalidGuild                                 JSONErrorCode = 50055
	JSONErrorCodeInvalidSKU                                   JSONErrorCode = 50057
	JSONErrorCodeInvalidRequestOrigin                         JSONErrorCode = 50067
	JSONErrorCodeInvalidMessageType                           JSONErrorCode = 50068
	JSONErrorCodePaymentSourceRequiredToRedeemGift            JSONErrorCode = 50070
	JSONErrorCodeCannotModifySystemWebhook                    JSONErrorCode = 50073
	JSONErrorCodeCannotDeleteChannelRequiredForCommunity      JSONErrorCode = 50074
	JSONErrorCodeCannotEditStickersWithinMessage              JSONErrorCode = 50080
	JSONErrorCodeInvalidStickerSent                           JSONErrorCode = 50081
	JSONErrorCodeOperationOnArchivedThread                    JSONErrorCode = 50083
	JSONErrorCodeInvalidThreadNotificationSettings            JSONErrorCode = 50084
	JSONErrorCodeBeforeValueEarlierThanThreadCreation         JSONErrorCode = 50085
	JSONErrorCodeCommunityServerChannelsMustBeText            JSONErrorCode = 50086
	JSONErrorCodeEntityTypeDifferentFromEventEntity           JSONErrorCode = 50091
	JSONErrorCodeServerNotAvailableInLocation                 JSONErrorCode = 50095
	JSONErrorCodeServerNeedsMonetizationEnabled               JSONErrorCode = 50097
	JSONErrorCodeServerNeedsMoreBoosts                        JSONErrorCode = 50101
	JSONErrorCodeRequestBodyContainsInvalidJSON               JSONErrorCode = 50109
	JSONErrorCodeProvidedFileInvalid                          JSONErrorCode = 50110
	JSONErrorCodeProvidedFileTypeInvalid                      JSONErrorCode = 50123
	JSONErrorCodeProvidedFileDurationExceedsMaximum           JSONErrorCode = 50124
	JSONErrorCodeOwnerCannotBePendingMember                   JSONErrorCode = 50131
	JSONErrorCodeOwnershipCannotBeTransferredToBot            JSONErrorCode = 50132
	JSONErrorCodeFailedToResizeAssetBelowMaximumSize          JSONErrorCode = 50138
	JSONErrorCodeCannotMixSubscriptionAndNonSubscriptionRoles JSONErrorCode = 50144
	JSONErrorCodeCannotConvertBetweenPremiumAndNormalEmoji    JSONErrorCode = 50145
	JSONErrorCodeUploadedFileNotFound                         JSONErrorCode = 50146
	JSONErrorCodeSpecifiedEmojiInvalid                        JSONErrorCode = 50151
	JSONErrorCodeVoiceMessagesDoNotSupportAdditionalContent   JSONErrorCode = 50159
	JSONErrorCodeVoiceMessagesMustHaveSingleAudioAttachment   JSONErrorCode = 50160
	JSONErrorCodeVoiceMessagesMustHaveSupportingMetadata      JSONErrorCode = 50161
	JSONErrorCodeVoiceMessagesCannotBeEdited                  JSONErrorCode = 50162
	JSONErrorCodeCannotDeleteGuildSubscriptionIntegration     JSONErrorCode = 50163
	JSONErrorCodeCannotSendVoiceMessagesInChannel             JSONErrorCode = 50173
	JSONErrorCodeUserAccountMustFirstBeVerified               JSONErrorCode = 50178
	JSONErrorCodeProvidedFileDoesNotHaveValidDuration         JSONErrorCode = 50192

	// Permission error
	JSONErrorCodeNoPermissionToSendSticker JSONErrorCode = 50600

	// Two factor required
	JSONErrorCodeTwoFactorRequired JSONErrorCode = 60003

	// No users with DiscordTag
	JSONErrorCodeNoUsersWithDiscordTagExist JSONErrorCode = 80004

	// Reaction errors
	JSONErrorCodeReactionBlocked             JSONErrorCode = 90001
	JSONErrorCodeUserCannotUseBurstReactions JSONErrorCode = 90002

	// Application not available
	JSONErrorCodeApplicationNotYetAvailable JSONErrorCode = 110001

	// API overloaded
	JSONErrorCodeAPIResourceOverloaded JSONErrorCode = 130000

	// Stage already open
	JSONErrorCodeStageAlreadyOpen JSONErrorCode = 150006

	// Thread errors
	JSONErrorCodeCannotReplyWithoutReadMessageHistoryPermission JSONErrorCode = 160002
	JSONErrorCodeThreadAlreadyCreatedForMessage                 JSONErrorCode = 160004
	JSONErrorCodeThreadLocked                                   JSONErrorCode = 160005
	JSONErrorCodeMaximumActiveThreadsReached                    JSONErrorCode = 160006
	JSONErrorCodeMaximumActiveAnnouncementThreadsReached        JSONErrorCode = 160007

	// Lottie/sticker errors
	JSONErrorCodeInvalidJSONForUploadedLottieFile             JSONErrorCode = 170001
	JSONErrorCodeUploadedLottiesCannotContainRasterizedImages JSONErrorCode = 170002
	JSONErrorCodeStickerMaximumFramerateExceeded              JSONErrorCode = 170003
	JSONErrorCodeStickerFrameCountExceedsMaximum              JSONErrorCode = 170004
	JSONErrorCodeLottieAnimationMaximumDimensionsExceeded     JSONErrorCode = 170005
	JSONErrorCodeStickerFrameRateTooSmallOrTooLarge           JSONErrorCode = 170006
	JSONErrorCodeStickerAnimationDurationExceedsMaximum       JSONErrorCode = 170007

	// Event errors
	JSONErrorCodeCannotUpdateFinishedEvent        JSONErrorCode = 180000
	JSONErrorCodeFailedToCreateStageForStageEvent JSONErrorCode = 180002

	// Auto moderation
	JSONErrorCodeMessageBlockedByAutomaticModeration JSONErrorCode = 200000
	JSONErrorCodeTitleBlockedByAutomaticModeration   JSONErrorCode = 200001

	// Webhook errors
	JSONErrorCodeWebhooksPostedToForumChannelsMustHaveThreadNameOrID        JSONErrorCode = 220001
	JSONErrorCodeWebhooksPostedToForumChannelsCannotHaveBothThreadNameAndID JSONErrorCode = 220002
	JSONErrorCodeWebhooksCanOnlyCreateThreadsInForumChannels                JSONErrorCode = 220003
	JSONErrorCodeWebhookServicesCannotBeUsedInForumChannels                 JSONErrorCode = 220004

	// Harmful links
	JSONErrorCodeMessageBlockedByHarmfulLinksFilter JSONErrorCode = 240000

	// Onboarding errors
	JSONErrorCodeCannotEnableOnboardingRequirementsNotMet     JSONErrorCode = 350000
	JSONErrorCodeCannotUpdateOnboardingWhileBelowRequirements JSONErrorCode = 350001

	// File upload limit
	JSONErrorCodeAccessToFileUploadsLimitedForGuild JSONErrorCode = 400001

	// Failed to ban users
	JSONErrorCodeFailedToBanUsers JSONErrorCode = 500000

	// Poll errors
	JSONErrorCodePollVotingBlocked                 JSONErrorCode = 520000
	JSONErrorCodePollExpired                       JSONErrorCode = 520001
	JSONErrorCodeInvalidChannelTypeForPollCreation JSONErrorCode = 520002
	JSONErrorCodeCannotEditPollMessage             JSONErrorCode = 520003
	JSONErrorCodeCannotUseEmojiIncludedWithPoll    JSONErrorCode = 520004
	JSONErrorCodeCannotExpireNonPollMessage        JSONErrorCode = 520006
)

var _ error = (*Error)(nil)

// Error holds the *[http.Request], *[http.Response] & an error related to a REST request.
// It's always a pointer to *[Error] that is returned by the REST client.
type Error struct {
	Request  *http.Request  `json:"-"`
	RqBody   []byte         `json:"-"`
	Response *http.Response `json:"-"`
	RsBody   []byte         `json:"-"`

	Code    JSONErrorCode   `json:"code"`
	Errors  json.RawMessage `json:"errors"`
	Message string          `json:"message"`
}

// newError returns a new *Error with the given http.Request, http.Response
func newError(rq *http.Request, rqBody []byte, rs *http.Response, rsBody []byte) *Error {
	err := &Error{
		Request:  rq,
		RqBody:   rqBody,
		Response: rs,
		RsBody:   rsBody,
	}
	_ = json.Unmarshal(rsBody, &err)

	return err
}

// Is returns true if the error is a *Error with the same status code as the target error
func (e *Error) Is(target error) bool {
	var err *Error
	if ok := errors.As(target, &err); !ok {
		return false
	}
	if e.Code != 0 && err.Code != 0 {
		return e.Code == err.Code
	}
	return err.Response != nil && e.Response != nil && err.Response.StatusCode == e.Response.StatusCode
}

// Error returns the error formatted as string
func (e *Error) Error() string {
	if e.Code != 0 {
		msg := fmt.Sprintf("%d: %s", e.Code, e.Message)
		if e.Code == 50035 {
			msg += fmt.Sprintf("\n%s", printErrors(e.Errors))
		}
		return msg
	}
	return fmt.Sprintf("Status: %s, Body: %s", e.Response.Status, string(e.RsBody))
}

// Error returns the error formatted as string
func (e *Error) String() string {
	return e.Error()
}

func printErrors(errors json.RawMessage) string {
	var m map[string]any
	if err := json.Unmarshal(errors, &m); err != nil {
		return ""
	}

	return parseErrors("", m)
}

func parseErrors(prefix string, err map[string]any) string {
	if errs, ok := err["_errors"]; ok {
		var s []string
		for _, e := range errs.([]any) {
			m := e.(map[string]any)
			s = append(s, fmt.Sprintf("%s -> %s: %s", prefix, m["code"], m["message"]))
		}
		return strings.Join(s, "\n")
	}

	var s []string
	for _, k := range slices.Sorted(maps.Keys(err)) {
		m := err[k].(map[string]any)

		nextPrefix := prefix
		if nextPrefix != "" {
			nextPrefix += " -> "
		}

		s = append(s, parseErrors(nextPrefix+k, m))
	}

	return strings.Join(s, "\n")
}

// IsJSONErrorCode returns true if the error is a *Error with one of the given error codes
func IsJSONErrorCode(err error, codes ...JSONErrorCode) bool {
	var restErr *Error
	if ok := errors.As(err, &restErr); !ok {
		return false
	}
	return slices.Contains(codes, restErr.Code)
}
