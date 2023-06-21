package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware/otelhandler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	Name       = "example"
	Namespace  = "github.com/disgoorg/disgo/_examples/otel"
	InstanceID = "1"
	Version    = "0.0.1"
)

var (
	token        = os.Getenv("disgo_token")
	guildID      = snowflake.GetEnv("disgo_guild_id")
	otelEndpoint = os.Getenv("otel_endpoint")
	otelSecure   = os.Getenv("otel_secure")
	commands     = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Replies with pong",
		},
	}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tracer, err := newTracer()
	if err != nil {
		log.Fatal("error while getting tracer")
	}

	r := handler.New()
	r.Use(otelhandler.Middleware("example"))
	r.Command("/ping", pingHandler(tracer))

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithRestClientConfigOpts(
			rest.WithHTTPClient(&http.Client{
				Transport: otelhttp.NewTransport(nil),
				Timeout:   5 * time.Second,
			}),
		),
		bot.WithEventListeners(r),
	)
	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	if err = handler.SyncCommands(client, commands, []snowflake.ID{guildID}); err != nil {
		log.Fatal("error while syncing commands: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func pingHandler(tracer trace.Tracer) func(event *handler.CommandEvent) error {
	return func(event *handler.CommandEvent) error {
		ctx, span := tracer.Start(event.Ctx, "ping",
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attribute.String("my.attribute", "test")),
		)
		defer span.End()

		return event.CreateMessage(discord.MessageCreate{
			Content: "pong",
		}, rest.WithCtx(ctx))
	}
}
