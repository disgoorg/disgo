# webhook

[Webhook](https://discord.com/developers/docs/resources/webhook) module of [disgo](https://github.com/disgoorg/disgo)

### Usage

Import the package into your project.

```go
import "github.com/disgoorg/disgo/webhook"
```

Create a new Webhook by `webhook_id` and `webhook_token`. (*This WebhookClient should be created once as it holds important state*)

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken")

client, err := webhook.NewWithURL("webhookURL")
```

`webhook.New` takes a vararg of type `webhook.ConfigOpt` as third argument which lets you pass additional optional parameter like a custom logger, rest client, etc

### Optional Arguments

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken",
	webhook.WithLogger(logrus.New()),
	webhook.WithDefaultAllowedMentions(discord.AllowedMentions{
		RepliedUser: false,
	}),
)
```

### Send Message

You can send a message as following

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken")

message, err := client.CreateContent("hello world!")

message, err := client.CreateEmbeds(discord.NewEmbedBuilder().
	SetDescription("hello world!").
	Build(),
)

message, err := client.CreateMessage(webhook.NewWebhookMessageCreateBuilder().
	SetContent("hello world!").
	Build(),
)

message, err := client.CreateMessage(discord.WebhookMessageCreate{
	Content: "hello world!",
})
```

### Edit Message

Messages can also be edited

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken")

message, err := client.UpdateContent("870741249114652722", "hello world!")

message, err := client.UpdateEmbeds("870741249114652722", discord.NewEmbedBuilder().
	SetDescription("hello world!").
	Build(),
)

message, err := client.UpdateMessage("870741249114652722", discord.NewWebhookMessageUpdateBuilder().
	SetContent("hello world!").
	Build(),
)

message, err := client.UpdateMessage("870741249114652722", discord.WebhookMessageUpdate{
	Content: json.Ptr("hello world!"),
})
```

### Delete Message

or deleted

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken")

err := client.DeleteMessage("message_id")
```

### Full Example

a full example can be found [here](https://github.com/disgoorg/disgo/tree/master/_examples/webhook/example.go)