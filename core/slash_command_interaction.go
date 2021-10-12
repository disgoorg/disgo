package core

import "github.com/DisgoOrg/disgo/discord"

type SlashCommandInteractionFilter func(slashCommandInteraction *SlashCommandInteraction) bool

type SlashCommandInteraction struct {
	discord.SlashCommandInteraction
	RespondInteraction
	CreateInteraction
	FollowupInteraction
	CommandID           discord.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Resolved            SlashCommandResolved
	Options             OptionsMap
}

// CommandPath returns the ApplicationCommand path
func (i *SlashCommandInteraction) CommandPath() string {
	path := i.Data.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

// Guild returns the Guild from the Caches
func (i *SlashCommandInteraction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.GuildCache().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *SlashCommandInteraction) Channel() *Channel {
	return i.Bot.Caches.ChannelCache().Get(i.ChannelID)
}

// SlashCommandResolved contains resolved mention data for SlashCommand(s)
type SlashCommandResolved struct {
	Users    map[discord.Snowflake]*User
	Members  map[discord.Snowflake]*Member
	Roles    map[discord.Snowflake]*Role
	Channels map[discord.Snowflake]*Channel
}

type OptionsMap map[string]SlashCommandOption

func (m OptionsMap) Get(name string) *SlashCommandOption {
	if option, ok := m[name]; ok {
		return &option
	}
	return nil
}

func (m OptionsMap) GetAll() []SlashCommandOption {
	options := make([]SlashCommandOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m OptionsMap) GetByType(optionType discord.ApplicationCommandOptionType) []SlashCommandOption {
	return m.FindAll(func(option SlashCommandOption) bool {
		return option.Type() == optionType
	})
}

func (m OptionsMap) Find(optionFindFunc func(option SlashCommandOption) bool) *SlashCommandOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return &option
		}
	}
	return nil
}

func (m OptionsMap) FindAll(optionFindFunc func(option SlashCommandOption) bool) []SlashCommandOption {
	var options []SlashCommandOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}
