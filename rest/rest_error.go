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

// ErrorCode is the error code returned by the Discord API.
// See https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
type ErrorCode int

const (
	// General error (such as a malformed request body, amongst other things)
	ErrorCodeGeneral ErrorCode = 0

	// Unknown resources
	ErrorCodeUnknownAccount                        ErrorCode = 10001
	ErrorCodeUnknownApplication                    ErrorCode = 10002
	ErrorCodeUnknownChannel                        ErrorCode = 10003
	ErrorCodeUnknownGuild                          ErrorCode = 10004
	ErrorCodeUnknownIntegration                    ErrorCode = 10005
	ErrorCodeUnknownInvite                         ErrorCode = 10006
	ErrorCodeUnknownMember                         ErrorCode = 10007
	ErrorCodeUnknownMessage                        ErrorCode = 10008
	ErrorCodeUnknownPermissionOverwrite            ErrorCode = 10009
	ErrorCodeUnknownProvider                       ErrorCode = 10010
	ErrorCodeUnknownRole                           ErrorCode = 10011
	ErrorCodeUnknownToken                          ErrorCode = 10012
	ErrorCodeUnknownUser                           ErrorCode = 10013
	ErrorCodeUnknownEmoji                          ErrorCode = 10014
	ErrorCodeUnknownWebhook                        ErrorCode = 10015
	ErrorCodeUnknownWebhookService                 ErrorCode = 10016
	ErrorCodeUnknownSession                        ErrorCode = 10020
	ErrorCodeUnknownAsset                          ErrorCode = 10021
	ErrorCodeUnknownBan                            ErrorCode = 10026
	ErrorCodeUnknownSKU                            ErrorCode = 10027
	ErrorCodeUnknownStoreListing                   ErrorCode = 10028
	ErrorCodeUnknownEntitlement                    ErrorCode = 10029
	ErrorCodeUnknownBuild                          ErrorCode = 10030
	ErrorCodeUnknownLobby                          ErrorCode = 10031
	ErrorCodeUnknownBranch                         ErrorCode = 10032
	ErrorCodeUnknownStoreDirectoryLayout           ErrorCode = 10033
	ErrorCodeUnknownRedistributable                ErrorCode = 10036
	ErrorCodeUnknownGiftCode                       ErrorCode = 10038
	ErrorCodeUnknownStream                         ErrorCode = 10049
	ErrorCodeUnknownPremiumServerSubscribeCooldown ErrorCode = 10050
	ErrorCodeUnknownGuildTemplate                  ErrorCode = 10057
	ErrorCodeUnknownDiscoverableServerCategory     ErrorCode = 10059
	ErrorCodeUnknownSticker                        ErrorCode = 10060
	ErrorCodeUnknownStickerPack                    ErrorCode = 10061
	ErrorCodeUnknownInteraction                    ErrorCode = 10062
	ErrorCodeUnknownApplicationCommand             ErrorCode = 10063
	ErrorCodeUnknownVoiceState                     ErrorCode = 10065
	ErrorCodeUnknownApplicationCommandPermissions  ErrorCode = 10066
	ErrorCodeUnknownStageInstance                  ErrorCode = 10067
	ErrorCodeUnknownGuildMemberVerificationForm    ErrorCode = 10068
	ErrorCodeUnknownGuildWelcomeScreen             ErrorCode = 10069
	ErrorCodeUnknownGuildScheduledEvent            ErrorCode = 10070
	ErrorCodeUnknownGuildScheduledEventUser        ErrorCode = 10071
	ErrorCodeUnknownTag                            ErrorCode = 10087
	ErrorCodeUnknownSound                          ErrorCode = 10097

	// Authorization/action errors
	ErrorCodeBotsCannotUseThisEndpoint                   ErrorCode = 20001
	ErrorCodeOnlyBotsCanUseThisEndpoint                  ErrorCode = 20002
	ErrorCodeExplicitContentCannotBeSent                 ErrorCode = 20009
	ErrorCodeNotAuthorizedToPerformThisAction            ErrorCode = 20012
	ErrorCodeActionCannotBePerformedDueToSlowmode        ErrorCode = 20016
	ErrorCodeOnlyOwnerCanPerformThisAction               ErrorCode = 20018
	ErrorCodeMessageCannotBeEditedDueToAnnouncementLimit ErrorCode = 20022
	ErrorCodeUnderMinimumAge                             ErrorCode = 20024
	ErrorCodeChannelWriteRateLimit                       ErrorCode = 20028
	ErrorCodeServerWriteRateLimit                        ErrorCode = 20029
	ErrorCodeStageTopicServerNameOrChannelNamesBlocked   ErrorCode = 20031
	ErrorCodeGuildPremiumSubscriptionLevelTooLow         ErrorCode = 20035

	// Maximum limits
	ErrorCodeMaximumGuildsReached                          ErrorCode = 30001
	ErrorCodeMaximumFriendsReached                         ErrorCode = 30002
	ErrorCodeMaximumPinsReached                            ErrorCode = 30003
	ErrorCodeMaximumRecipientsReached                      ErrorCode = 30004
	ErrorCodeMaximumGuildRolesReached                      ErrorCode = 30005
	ErrorCodeMaximumWebhooksReached                        ErrorCode = 30007
	ErrorCodeMaximumEmojisReached                          ErrorCode = 30008
	ErrorCodeMaximumReactionsReached                       ErrorCode = 30010
	ErrorCodeMaximumGroupDMsReached                        ErrorCode = 30011
	ErrorCodeMaximumGuildChannelsReached                   ErrorCode = 30013
	ErrorCodeMaximumAttachmentsInMessageReached            ErrorCode = 30015
	ErrorCodeMaximumInvitesReached                         ErrorCode = 30016
	ErrorCodeMaximumAnimatedEmojisReached                  ErrorCode = 30018
	ErrorCodeMaximumServerMembersReached                   ErrorCode = 30019
	ErrorCodeMaximumServerCategoriesReached                ErrorCode = 30030
	ErrorCodeGuildAlreadyHasTemplate                       ErrorCode = 30031
	ErrorCodeMaximumApplicationCommandsReached             ErrorCode = 30032
	ErrorCodeMaximumThreadParticipantsReached              ErrorCode = 30033
	ErrorCodeMaximumDailyApplicationCommandCreatesReached  ErrorCode = 30034
	ErrorCodeMaximumBansForNonGuildMembersExceeded         ErrorCode = 30035
	ErrorCodeMaximumBansFetchesReached                     ErrorCode = 30037
	ErrorCodeMaximumUncompletedGuildScheduledEventsReached ErrorCode = 30038
	ErrorCodeMaximumStickersReached                        ErrorCode = 30039
	ErrorCodeMaximumPruneRequestsReached                   ErrorCode = 30040
	ErrorCodeMaximumGuildWidgetSettingsUpdatesReached      ErrorCode = 30042
	ErrorCodeMaximumSoundboardSoundsReached                ErrorCode = 30045
	ErrorCodeMaximumEditsToMessagesOlderThan1HourReached   ErrorCode = 30046
	ErrorCodeMaximumPinnedThreadsInForumChannelReached     ErrorCode = 30047
	ErrorCodeMaximumTagsInForumChannelReached              ErrorCode = 30048
	ErrorCodeBitrateTooHighForChannelType                  ErrorCode = 30052
	ErrorCodeMaximumPremiumEmojisReached                   ErrorCode = 30056
	ErrorCodeMaximumWebhooksPerGuildReached                ErrorCode = 30058
	ErrorCodeMaximumChannelPermissionOverwritesReached     ErrorCode = 30060
	ErrorCodeChannelsForGuildTooLarge                      ErrorCode = 30061

	// Temporary/rate limit errors
	ErrorCodeUnauthorized                            ErrorCode = 40001
	ErrorCodeNeedToVerifyAccount                     ErrorCode = 40002
	ErrorCodeOpeningDirectMessagesTooFast            ErrorCode = 40003
	ErrorCodeSendMessagesTemporarilyDisabled         ErrorCode = 40004
	ErrorCodeRequestEntityTooLarge                   ErrorCode = 40005
	ErrorCodeFeatureTemporarilyDisabled              ErrorCode = 40006
	ErrorCodeUserBannedFromGuild                     ErrorCode = 40007
	ErrorCodeConnectionRevoked                       ErrorCode = 40012
	ErrorCodeOnlyConsumableSKUsCanBeConsumed         ErrorCode = 40018
	ErrorCodeCanOnlyDeleteSandboxEntitlements        ErrorCode = 40019
	ErrorCodeTargetUserNotConnectedToVoice           ErrorCode = 40032
	ErrorCodeMessageAlreadyCrossposted               ErrorCode = 40033
	ErrorCodeApplicationCommandWithNameAlreadyExists ErrorCode = 40041
	ErrorCodeApplicationInteractionFailedToSend      ErrorCode = 40043
	ErrorCodeCannotSendMessageInForumChannel         ErrorCode = 40058
	ErrorCodeInteractionAlreadyAcknowledged          ErrorCode = 40060
	ErrorCodeTagNamesMustBeUnique                    ErrorCode = 40061
	ErrorCodeServiceResourceRateLimited              ErrorCode = 40062
	ErrorCodeNoTagsAvailableForNonModerators         ErrorCode = 40066
	ErrorCodeTagRequiredToCreateForumPost            ErrorCode = 40067
	ErrorCodeEntitlementAlreadyGrantedForResource    ErrorCode = 40074
	ErrorCodeInteractionMaxFollowUpMessagesReached   ErrorCode = 40094

	// Cloudflare blocking
	ErrorCodeCloudflareBlockingRequest ErrorCode = 40333

	// Invalid/missing errors
	ErrorCodeMissingAccess                                ErrorCode = 50001
	ErrorCodeInvalidAccountType                           ErrorCode = 50002
	ErrorCodeCannotExecuteActionOnDMChannel               ErrorCode = 50003
	ErrorCodeGuildWidgetDisabled                          ErrorCode = 50004
	ErrorCodeCannotEditMessageAuthoredByAnotherUser       ErrorCode = 50005
	ErrorCodeCannotSendEmptyMessage                       ErrorCode = 50006
	ErrorCodeCannotSendMessagesToThisUser                 ErrorCode = 50007
	ErrorCodeCannotSendMessagesInNonTextChannel           ErrorCode = 50008
	ErrorCodeChannelVerificationLevelTooHigh              ErrorCode = 50009
	ErrorCodeOAuth2ApplicationDoesNotHaveBot              ErrorCode = 50010
	ErrorCodeOAuth2ApplicationLimitReached                ErrorCode = 50011
	ErrorCodeInvalidOAuth2State                           ErrorCode = 50012
	ErrorCodeLackPermissionsToPerformAction               ErrorCode = 50013
	ErrorCodeInvalidAuthenticationToken                   ErrorCode = 50014
	ErrorCodeNoteTooLong                                  ErrorCode = 50015
	ErrorCodeTooFewOrTooManyMessagesToDelete              ErrorCode = 50016
	ErrorCodeInvalidMFALevel                              ErrorCode = 50017
	ErrorCodeMessageCanOnlyBePinnedToOriginalChannel      ErrorCode = 50019
	ErrorCodeInviteCodeInvalidOrTaken                     ErrorCode = 50020
	ErrorCodeCannotExecuteActionOnSystemMessage           ErrorCode = 50021
	ErrorCodeCannotExecuteActionOnThisChannelType         ErrorCode = 50024
	ErrorCodeInvalidOAuth2AccessToken                     ErrorCode = 50025
	ErrorCodeMissingRequiredOAuth2Scope                   ErrorCode = 50026
	ErrorCodeInvalidWebhookToken                          ErrorCode = 50027
	ErrorCodeInvalidRole                                  ErrorCode = 50028
	ErrorCodeInvalidRecipients                            ErrorCode = 50033
	ErrorCodeMessageTooOldToBulkDelete                    ErrorCode = 50034
	ErrorCodeInvalidFormBody                              ErrorCode = 50035
	ErrorCodeInviteAcceptedToGuildBotNotIn                ErrorCode = 50036
	ErrorCodeInvalidActivityAction                        ErrorCode = 50039
	ErrorCodeInvalidAPIVersion                            ErrorCode = 50041
	ErrorCodeFileUploadedExceedsMaximumSize               ErrorCode = 50045
	ErrorCodeInvalidFileUploaded                          ErrorCode = 50046
	ErrorCodeCannotSelfRedeemGift                         ErrorCode = 50054
	ErrorCodeInvalidGuild                                 ErrorCode = 50055
	ErrorCodeInvalidSKU                                   ErrorCode = 50057
	ErrorCodeInvalidRequestOrigin                         ErrorCode = 50067
	ErrorCodeInvalidMessageType                           ErrorCode = 50068
	ErrorCodePaymentSourceRequiredToRedeemGift            ErrorCode = 50070
	ErrorCodeCannotModifySystemWebhook                    ErrorCode = 50073
	ErrorCodeCannotDeleteChannelRequiredForCommunity      ErrorCode = 50074
	ErrorCodeCannotEditStickersWithinMessage              ErrorCode = 50080
	ErrorCodeInvalidStickerSent                           ErrorCode = 50081
	ErrorCodeOperationOnArchivedThread                    ErrorCode = 50083
	ErrorCodeInvalidThreadNotificationSettings            ErrorCode = 50084
	ErrorCodeBeforeValueEarlierThanThreadCreation         ErrorCode = 50085
	ErrorCodeCommunityServerChannelsMustBeText            ErrorCode = 50086
	ErrorCodeEntityTypeDifferentFromEventEntity           ErrorCode = 50091
	ErrorCodeServerNotAvailableInLocation                 ErrorCode = 50095
	ErrorCodeServerNeedsMonetizationEnabled               ErrorCode = 50097
	ErrorCodeServerNeedsMoreBoosts                        ErrorCode = 50101
	ErrorCodeRequestBodyContainsInvalidJSON               ErrorCode = 50109
	ErrorCodeProvidedFileInvalid                          ErrorCode = 50110
	ErrorCodeProvidedFileTypeInvalid                      ErrorCode = 50123
	ErrorCodeProvidedFileDurationExceedsMaximum           ErrorCode = 50124
	ErrorCodeOwnerCannotBePendingMember                   ErrorCode = 50131
	ErrorCodeOwnershipCannotBeTransferredToBot            ErrorCode = 50132
	ErrorCodeFailedToResizeAssetBelowMaximumSize          ErrorCode = 50138
	ErrorCodeCannotMixSubscriptionAndNonSubscriptionRoles ErrorCode = 50144
	ErrorCodeCannotConvertBetweenPremiumAndNormalEmoji    ErrorCode = 50145
	ErrorCodeUploadedFileNotFound                         ErrorCode = 50146
	ErrorCodeSpecifiedEmojiInvalid                        ErrorCode = 50151
	ErrorCodeVoiceMessagesDoNotSupportAdditionalContent   ErrorCode = 50159
	ErrorCodeVoiceMessagesMustHaveSingleAudioAttachment   ErrorCode = 50160
	ErrorCodeVoiceMessagesMustHaveSupportingMetadata      ErrorCode = 50161
	ErrorCodeVoiceMessagesCannotBeEdited                  ErrorCode = 50162
	ErrorCodeCannotDeleteGuildSubscriptionIntegration     ErrorCode = 50163
	ErrorCodeCannotSendVoiceMessagesInChannel             ErrorCode = 50173
	ErrorCodeUserAccountMustFirstBeVerified               ErrorCode = 50178
	ErrorCodeProvidedFileDoesNotHaveValidDuration         ErrorCode = 50192

	// Permission error
	ErrorCodeNoPermissionToSendSticker ErrorCode = 50600

	// Two factor required
	ErrorCodeTwoFactorRequired ErrorCode = 60003

	// No users with DiscordTag
	ErrorCodeNoUsersWithDiscordTagExist ErrorCode = 80004

	// Reaction errors
	ErrorCodeReactionBlocked             ErrorCode = 90001
	ErrorCodeUserCannotUseBurstReactions ErrorCode = 90002

	// Application not available
	ErrorCodeApplicationNotYetAvailable ErrorCode = 110001

	// API overloaded
	ErrorCodeAPIResourceOverloaded ErrorCode = 130000

	// Stage already open
	ErrorCodeStageAlreadyOpen ErrorCode = 150006

	// Thread errors
	ErrorCodeCannotReplyWithoutReadMessageHistoryPermission ErrorCode = 160002
	ErrorCodeThreadAlreadyCreatedForMessage                 ErrorCode = 160004
	ErrorCodeThreadLocked                                   ErrorCode = 160005
	ErrorCodeMaximumActiveThreadsReached                    ErrorCode = 160006
	ErrorCodeMaximumActiveAnnouncementThreadsReached        ErrorCode = 160007

	// Lottie/sticker errors
	ErrorCodeInvalidJSONForUploadedLottieFile             ErrorCode = 170001
	ErrorCodeUploadedLottiesCannotContainRasterizedImages ErrorCode = 170002
	ErrorCodeStickerMaximumFramerateExceeded              ErrorCode = 170003
	ErrorCodeStickerFrameCountExceedsMaximum              ErrorCode = 170004
	ErrorCodeLottieAnimationMaximumDimensionsExceeded     ErrorCode = 170005
	ErrorCodeStickerFrameRateTooSmallOrTooLarge           ErrorCode = 170006
	ErrorCodeStickerAnimationDurationExceedsMaximum       ErrorCode = 170007

	// Event errors
	ErrorCodeCannotUpdateFinishedEvent        ErrorCode = 180000
	ErrorCodeFailedToCreateStageForStageEvent ErrorCode = 180002

	// Auto moderation
	ErrorCodeMessageBlockedByAutomaticModeration ErrorCode = 200000
	ErrorCodeTitleBlockedByAutomaticModeration   ErrorCode = 200001

	// Webhook errors
	ErrorCodeWebhooksPostedToForumChannelsMustHaveThreadNameOrID        ErrorCode = 220001
	ErrorCodeWebhooksPostedToForumChannelsCannotHaveBothThreadNameAndID ErrorCode = 220002
	ErrorCodeWebhooksCanOnlyCreateThreadsInForumChannels                ErrorCode = 220003
	ErrorCodeWebhookServicesCannotBeUsedInForumChannels                 ErrorCode = 220004

	// Harmful links
	ErrorCodeMessageBlockedByHarmfulLinksFilter ErrorCode = 240000

	// Onboarding errors
	ErrorCodeCannotEnableOnboardingRequirementsNotMet     ErrorCode = 350000
	ErrorCodeCannotUpdateOnboardingWhileBelowRequirements ErrorCode = 350001

	// File upload limit
	ErrorCodeAccessToFileUploadsLimitedForGuild ErrorCode = 400001

	// Failed to ban users
	ErrorCodeFailedToBanUsers ErrorCode = 500000

	// Poll errors
	ErrorCodePollVotingBlocked                 ErrorCode = 520000
	ErrorCodePollExpired                       ErrorCode = 520001
	ErrorCodeInvalidChannelTypeForPollCreation ErrorCode = 520002
	ErrorCodeCannotEditPollMessage             ErrorCode = 520003
	ErrorCodeCannotUseEmojiIncludedWithPoll    ErrorCode = 520004
	ErrorCodeCannotExpireNonPollMessage        ErrorCode = 520006
)

var _ error = (*Error)(nil)

// Error holds the *[http.Request], *[http.Response] & an error related to a REST request.
// It's always a pointer to *[Error] that is returned by the REST client.
type Error struct {
	Request  *http.Request  `json:"-"`
	RqBody   []byte         `json:"-"`
	Response *http.Response `json:"-"`
	RsBody   []byte         `json:"-"`

	Code    ErrorCode       `json:"code"`
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

// IsRestErrorCode returns true if the error is a *Error with one of the given error codes
func IsRestErrorCode(err error, codes ...ErrorCode) bool {
	var restErr *Error
	if ok := errors.As(err, &restErr); !ok {
		return false
	}
	return slices.Contains(codes, restErr.Code)
}
