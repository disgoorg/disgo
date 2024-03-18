package handler

import (
	"os"
	"testing"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/stretchr/testify/assert"
)

func NewRecorder() *InteractionResponseRecorder {
	return &InteractionResponseRecorder{}
}

type InteractionResponseRecorder struct {
	Response *discord.InteractionResponse
}

func (i *InteractionResponseRecorder) Respond(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
	i.Response = &discord.InteractionResponse{
		Type: responseType,
		Data: data,
	}
	return nil
}

func TestCommandMux(t *testing.T) {
	slashData, err := os.ReadFile("testdata/command/slash_command.json")
	assert.NoError(t, err)

	userData, err := os.ReadFile("testdata/command/user_command.json")
	assert.NoError(t, err)

	messageData, err := os.ReadFile("testdata/command/message_command.json")
	assert.NoError(t, err)

	data := []struct {
		data     []byte
		expected *discord.InteractionResponse
	}{
		{
			data: slashData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar",
				},
			},
		},
		{
			data: userData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar2",
				},
			},
		},
		{
			data: messageData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar3",
				},
			},
		},
	}

	mux := New()
	mux.SlashCommand("/foo", func(data discord.SlashCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "bar",
		})
	})
	mux.UserCommand("/foo", func(data discord.UserCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "bar2",
		})
	})
	mux.MessageCommand("/foo", func(data discord.MessageCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "bar3",
		})
	})

	for _, d := range data {
		interaction, err := discord.UnmarshalInteraction(d.data)
		assert.NoError(t, err)

		recorder := NewRecorder()
		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(nil, 0, 0),
			Interaction:  interaction,
			Respond:      recorder.Respond,
		})
		assert.Equal(t, d.expected, recorder.Response)
	}
}

func TestComponentMux(t *testing.T) {
	buttonData, err := os.ReadFile("testdata/component/button_component.json")
	assert.NoError(t, err)

	selectMenuData, err := os.ReadFile("testdata/component/select_menu_component.json")
	assert.NoError(t, err)

	data := []struct {
		data     []byte
		expected *discord.InteractionResponse
	}{
		{
			data: buttonData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar",
				},
			},
		},
		{
			data: selectMenuData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar2",
				},
			},
		},
	}

	mux := New()
	mux.ButtonComponent("/foo", func(data discord.ButtonInteractionData, e *ComponentEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "bar",
		})
	})
	mux.SelectMenuComponent("/foo", func(data discord.SelectMenuInteractionData, e *ComponentEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "bar2",
		})
	})

	for _, d := range data {
		interaction, err := discord.UnmarshalInteraction(d.data)
		assert.NoError(t, err)

		recorder := NewRecorder()
		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(nil, 0, 0),
			Interaction:  interaction,
			Respond:      recorder.Respond,
		})
		assert.Equal(t, d.expected, recorder.Response)
	}
}
