package handler

import (
	"maps"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
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

	entryPointData, err := os.ReadFile("testdata/command/entry_point_command.json")
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
		{
			data: entryPointData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar4",
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
	mux.EntryPointCommand("/foo", func(data discord.EntryPointCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "bar4",
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

func TestMiddlewareMux(t *testing.T) {
	slashData, err := os.ReadFile("testdata/middleware/slash_command.json")
	assert.NoError(t, err)

	data := []struct {
		data          []byte
		expected      *discord.InteractionResponse
		expectedVars1 map[string]string
		expectedVars2 map[string]string
		expectedVars3 map[string]string
	}{
		{
			data: slashData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "bar",
				},
			},
			expectedVars1: map[string]string{},
			expectedVars2: map[string]string{"bar": "bar"},
			expectedVars3: map[string]string{"bar": "bar"},
		},
	}

	var (
		dataVars1 = make(map[string]string)
		dataVars2 = make(map[string]string)
		dataVars3 = make(map[string]string)
	)

	mux := New()
	mux.Use(func(next Handler) Handler {
		return func(e *InteractionEvent) error {
			maps.Copy(dataVars1, e.Vars)
			return next(e)
		}
	})
	mux.Route("/foo/{bar}", func(r Router) {
		r.Use(func(next Handler) Handler {
			return func(e *InteractionEvent) error {
				maps.Copy(dataVars2, e.Vars)
				return next(e)
			}
		})
		r.Command("/baz", func(e *CommandEvent) error {
			maps.Copy(dataVars3, e.Vars)
			return e.CreateMessage(discord.MessageCreate{
				Content: "bar",
			})
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
		assert.Equal(t, d.expectedVars1, dataVars1)
		assert.Equal(t, d.expectedVars2, dataVars2)
		assert.Equal(t, d.expectedVars3, dataVars3)
	}
}

func TestMux(t *testing.T) {
	buttonFooData, err := os.ReadFile("testdata/mux/button_foo_component.json")
	assert.NoError(t, err)

	buttonFooBarData, err := os.ReadFile("testdata/mux/button_foo_bar_component.json")
	assert.NoError(t, err)

	data := []struct {
		data     []byte
		expected *discord.InteractionResponse
	}{
		{
			data: buttonFooData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "/foo/",
				},
			},
		},
		{
			data: buttonFooBarData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{
					Content: "/foo/bar",
				},
			},
		},
	}

	mux := New()
	mux.Route("/foo", func(r Router) {
		r.ButtonComponent("/", func(data discord.ButtonInteractionData, e *ComponentEvent) error {
			return e.CreateMessage(discord.MessageCreate{
				Content: "/foo/",
			})
		})
		r.ButtonComponent("/bar", func(data discord.ButtonInteractionData, e *ComponentEvent) error {
			return e.CreateMessage(discord.MessageCreate{
				Content: "/foo/bar",
			})
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
