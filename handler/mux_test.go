package handler

import (
	"context"
	"maps"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/disgoorg/disgo/bot"
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
	if err != nil {
		t.Fatalf("failed to read slash command data: %v", err)
	}

	userData, err := os.ReadFile("testdata/command/user_command.json")
	if err != nil {
		t.Fatalf("failed to read user command data: %v", err)
	}

	messageData, err := os.ReadFile("testdata/command/message_command.json")
	if err != nil {
		t.Fatalf("failed to read message command data: %v", err)
	}

	entryPointData, err := os.ReadFile("testdata/command/entry_point_command.json")
	if err != nil {
		t.Fatalf("failed to read entry point command data: %v", err)
	}

	data := []struct {
		data     []byte
		expected *discord.InteractionResponse
	}{
		{
			data: slashData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "bar"},
			},
		},
		{
			data: userData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "bar2"},
			},
		},
		{
			data: messageData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "bar3"},
			},
		},
		{
			data: entryPointData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "bar4"},
			},
		},
	}

	mux := New()
	mux.SlashCommand("/foo", func(data discord.SlashCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{Content: "bar"})
	})
	mux.UserCommand("/foo", func(data discord.UserCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{Content: "bar2"})
	})
	mux.MessageCommand("/foo", func(data discord.MessageCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{Content: "bar3"})
	})
	mux.EntryPointCommand("/foo", func(data discord.EntryPointCommandInteractionData, e *CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{Content: "bar4"})
	})

	for _, d := range data {
		interaction, err := discord.UnmarshalInteraction(d.data)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		recorder := NewRecorder()
		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(nil, 0, 0),
			Interaction:  interaction,
			Respond:      recorder.Respond,
		})

		if !reflect.DeepEqual(d.expected, recorder.Response) {
			t.Errorf("expected %+v, got %+v", d.expected, recorder.Response)
		}
	}
}

func TestComponentMux(t *testing.T) {
	buttonData, err := os.ReadFile("testdata/component/button_component.json")
	if err != nil {
		t.Fatalf("failed to read button component data: %v", err)
	}

	selectMenuData, err := os.ReadFile("testdata/component/select_menu_component.json")
	if err != nil {
		t.Fatalf("failed to read select menu component data: %v", err)
	}

	data := []struct {
		data     []byte
		expected *discord.InteractionResponse
	}{
		{
			data: buttonData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "bar"},
			},
		},
		{
			data: selectMenuData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "bar2"},
			},
		},
	}

	mux := New()
	mux.ButtonComponent("/foo", func(data discord.ButtonInteractionData, e *ComponentEvent) error {
		return e.CreateMessage(discord.MessageCreate{Content: "bar"})
	})
	mux.SelectMenuComponent("/foo", func(data discord.SelectMenuInteractionData, e *ComponentEvent) error {
		return e.CreateMessage(discord.MessageCreate{Content: "bar2"})
	})

	for _, d := range data {
		interaction, err := discord.UnmarshalInteraction(d.data)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		recorder := NewRecorder()
		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(nil, 0, 0),
			Interaction:  interaction,
			Respond:      recorder.Respond,
		})

		if !reflect.DeepEqual(d.expected, recorder.Response) {
			t.Errorf("expected %+v, got %+v", d.expected, recorder.Response)
		}
	}
}

func TestMiddlewareMux(t *testing.T) {
	slashData, err := os.ReadFile("testdata/middleware/slash_command.json")
	if err != nil {
		t.Fatalf("failed to read slash command data: %v", err)
	}

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
				Data: discord.MessageCreate{Content: "bar"},
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
			return e.CreateMessage(discord.MessageCreate{Content: "bar"})
		})
	})

	for _, d := range data {
		interaction, err := discord.UnmarshalInteraction(d.data)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		recorder := NewRecorder()
		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(nil, 0, 0),
			Interaction:  interaction,
			Respond:      recorder.Respond,
		})

		if !reflect.DeepEqual(d.expected, recorder.Response) {
			t.Errorf("expected %+v, got %+v", d.expected, recorder.Response)
		}
		if !reflect.DeepEqual(d.expectedVars1, dataVars1) {
			t.Errorf("expected vars1 %+v, got %+v", d.expectedVars1, dataVars1)
		}
		if !reflect.DeepEqual(d.expectedVars2, dataVars2) {
			t.Errorf("expected vars2 %+v, got %+v", d.expectedVars2, dataVars2)
		}
		if !reflect.DeepEqual(d.expectedVars3, dataVars3) {
			t.Errorf("expected vars3 %+v, got %+v", d.expectedVars3, dataVars3)
		}
	}
}

func TestMux(t *testing.T) {
	buttonFooData, err := os.ReadFile("testdata/mux/button_foo_component.json")
	if err != nil {
		t.Fatalf("failed to read button foo data: %v", err)
	}

	buttonFooBarData, err := os.ReadFile("testdata/mux/button_foo_bar_component.json")
	if err != nil {
		t.Fatalf("failed to read button foo bar data: %v", err)
	}

	data := []struct {
		data     []byte
		expected *discord.InteractionResponse
	}{
		{
			data: buttonFooData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "/foo/"},
			},
		},
		{
			data: buttonFooBarData,
			expected: &discord.InteractionResponse{
				Type: discord.InteractionResponseTypeCreateMessage,
				Data: discord.MessageCreate{Content: "/foo/bar"},
			},
		},
	}

	mux := New()
	mux.Route("/foo", func(r Router) {
		r.ButtonComponent("/", func(data discord.ButtonInteractionData, e *ComponentEvent) error {
			return e.CreateMessage(discord.MessageCreate{Content: "/foo/"})
		})
		r.ButtonComponent("/bar", func(data discord.ButtonInteractionData, e *ComponentEvent) error {
			return e.CreateMessage(discord.MessageCreate{Content: "/foo/bar"})
		})
	})

	for _, d := range data {
		interaction, err := discord.UnmarshalInteraction(d.data)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		recorder := NewRecorder()
		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(nil, 0, 0),
			Interaction:  interaction,
			Respond:      recorder.Respond,
		})

		if !reflect.DeepEqual(d.expected, recorder.Response) {
			t.Errorf("expected %+v, got %+v", d.expected, recorder.Response)
		}
	}
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestContextPropagation(t *testing.T) {
	slashData, err := os.ReadFile("testdata/command/slash_command.json")
	if err != nil {
		t.Fatalf("failed to read slash command data: %v", err)
	}

	type ctxKey struct{}

	var capturedCtxs []context.Context
	restClient := rest.NewClient("",
		rest.WithRateLimiter(rest.NewNoopRateLimiter()),
		rest.WithHTTPClient(&http.Client{
			Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				capturedCtxs = append(capturedCtxs, req.Context())
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       http.NoBody,
				}, nil
			}),
		}),
	)
	client := &bot.Client{
		Rest: rest.New(restClient),
	}

	t.Run("default event context is propagated", func(t *testing.T) {
		capturedCtxs = nil

		mux := New()
		mux.Use(func(next Handler) Handler {
			return func(e *InteractionEvent) error {
				e.Ctx = context.WithValue(e.Ctx, ctxKey{}, "middleware-ctx")
				return next(e)
			}
		})
		mux.SlashCommand("/foo", func(data discord.SlashCommandInteractionData, e *CommandEvent) error {
			if err := e.CreateMessage(discord.MessageCreate{Content: "bar"}); err != nil {
				return err
			}
			return e.DeleteInteractionResponse()
		})

		interaction, err := discord.UnmarshalInteraction(slashData)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(client, 0, 0),
			Interaction:  interaction,
			Respond: func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
				return restClient.Do(rest.CreateInteractionResponse.Compile(nil, interaction.ID(), interaction.Token()), nil, nil, opts...)
			},
		})

		if len(capturedCtxs) != 2 {
			t.Fatalf("expected 2 requests, got %d", len(capturedCtxs))
		}
		for i, ctx := range capturedCtxs {
			if got := ctx.Value(ctxKey{}); got != "middleware-ctx" {
				t.Errorf("request %d: expected context value %q, got %v", i, "middleware-ctx", got)
			}
		}
	})

	t.Run("reassigned handler event context is propagated", func(t *testing.T) {
		capturedCtxs = nil

		mux := New()
		mux.Use(func(next Handler) Handler {
			return func(e *InteractionEvent) error {
				e.Ctx = context.WithValue(e.Ctx, ctxKey{}, "middleware-ctx")
				return next(e)
			}
		})
		mux.SlashCommand("/foo", func(data discord.SlashCommandInteractionData, e *CommandEvent) error {
			e.Ctx = context.WithValue(e.Ctx, ctxKey{}, "handler-ctx")
			if err := e.CreateMessage(discord.MessageCreate{Content: "bar"}); err != nil {
				return err
			}
			return e.DeleteInteractionResponse()
		})

		interaction, err := discord.UnmarshalInteraction(slashData)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(client, 0, 0),
			Interaction:  interaction,
			Respond: func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
				return restClient.Do(rest.CreateInteractionResponse.Compile(nil, interaction.ID(), interaction.Token()), nil, nil, opts...)
			},
		})

		if len(capturedCtxs) != 2 {
			t.Fatalf("expected 2 requests, got %d", len(capturedCtxs))
		}
		for i, ctx := range capturedCtxs {
			if got := ctx.Value(ctxKey{}); got != "handler-ctx" {
				t.Errorf("request %d: expected context value %q, got %v", i, "handler-ctx", got)
			}
		}
	})

	t.Run("explicit WithCtx overrides event context", func(t *testing.T) {
		capturedCtxs = nil

		mux := New()
		mux.Use(func(next Handler) Handler {
			return func(e *InteractionEvent) error {
				e.Ctx = context.WithValue(e.Ctx, ctxKey{}, "middleware-ctx")
				return next(e)
			}
		})
		mux.SlashCommand("/foo", func(data discord.SlashCommandInteractionData, e *CommandEvent) error {
			customCtx := context.WithValue(context.Background(), ctxKey{}, "custom-ctx")
			if err := e.CreateMessage(discord.MessageCreate{Content: "bar"}, rest.WithCtx(customCtx)); err != nil {
				return err
			}
			return e.DeleteInteractionResponse(rest.WithCtx(customCtx))
		})

		interaction, err := discord.UnmarshalInteraction(slashData)
		if err != nil {
			t.Fatalf("failed to unmarshal interaction: %v", err)
		}

		mux.OnEvent(&events.InteractionCreate{
			GenericEvent: events.NewGenericEvent(client, 0, 0),
			Interaction:  interaction,
			Respond: func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
				return restClient.Do(rest.CreateInteractionResponse.Compile(nil, interaction.ID(), interaction.Token()), nil, nil, opts...)
			},
		})

		if len(capturedCtxs) != 2 {
			t.Fatalf("expected 2 requests, got %d", len(capturedCtxs))
		}
		for i, ctx := range capturedCtxs {
			if got := ctx.Value(ctxKey{}); got != "custom-ctx" {
				t.Errorf("request %d: expected context value %q, got %v", i, "custom-ctx", got)
			}
		}
	})
}
