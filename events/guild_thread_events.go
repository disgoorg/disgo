package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericThread is the base struct for all Thread events.
type GenericThread struct {
	*GenericEvent
	Thread   discord.GuildThread
	ThreadID snowflake.ID
	GuildID  snowflake.ID
	ParentID snowflake.ID
}

// ThreadCreate is dispatched when a thread is created.
type ThreadCreate struct {
	*GenericThread
	ThreadMember discord.ThreadMember
}

// ThreadUpdate is dispatched when a thread is updated.
type ThreadUpdate struct {
	*GenericThread
	OldThread discord.GuildThread
}

// ThreadDelete is dispatched when a thread is deleted.
type ThreadDelete struct {
	*GenericThread
}

// ThreadShow is dispatched when your bot gains access to a thread.
type ThreadShow struct {
	*GenericThread
}

// ThreadHide is dispatched when your bot loses access to a thread.
type ThreadHide struct {
	*GenericThread
}

// GenericThreadMember is the base struct for all ThreadMember events.
type GenericThreadMember struct {
	*GenericEvent
	GuildID        snowflake.ID
	ThreadID       snowflake.ID
	ThreadMemberID snowflake.ID
	ThreadMember   discord.ThreadMember
}

// ThreadMemberAdd is dispatched when a user is added to a thread.
type ThreadMemberAdd struct {
	*GenericThreadMember
	Member   discord.Member
	Presence *discord.Presence
}

// ThreadMemberUpdate is dispatched when a user is updated in a thread.
type ThreadMemberUpdate struct {
	*GenericThreadMember
	OldThreadMember discord.ThreadMember
}

// ThreadMemberRemove is dispatched when a user is removed from a thread.
type ThreadMemberRemove struct {
	*GenericThreadMember
}
