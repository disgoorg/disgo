package discord

// MessageActivityType is the type of MessageActivity https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
type MessageActivityType int

//Constants for MessageActivityType
//goland:noinspection GoUnusedConst
const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	_
	MessageActivityTypeJoinRequest
)
