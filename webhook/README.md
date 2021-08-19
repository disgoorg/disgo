# webhook

[WebHook](https://discord.com/developers/docs/resources/webhook) module of [disgo](https://github.com/DisgoOrg/disgo)

### Usage

Import the package into your project.

```go
import "github.com/DisgoOrg/disgo/webhook"
```

Create a new Webhook by `webhook_id/webhook_token`. (*This WebhookClient should be created once as it holds important
state*)

As first param you can optionally pass your own [*http.Client](https://pkg.go.dev/net/http#Client), as second you can
pass your own [rest.HTTPClient](https://pkg.go.dev/github.com/DisgoOrg/disgo/rest#HTTPClient)
and as third parameter you can pass your own logger
implementing [this](https://github.com/DisgoOrg/log/blob/master/logger.go) interface. This webhook then can be used to
send, edit and delete messages

### Send Message

```go
client, err := webhook.New(nil, nil, nil, "webhook_id", "webhook_token")

message, err := client.CreateContent(context.ToDo(), "hello world!")
message, err := client.CreateEmbeds(context.ToDo(), webhook.NewEmbedBuilder().
	SetDescription("hello world!").
	Build(),
)
message, err := client.CreateMessage(context.ToDo(), webhook.NewMessageCreateBuilder().
	SetContent("hello world!").
	Build(),
)
```

### Edit Message

```go
client, err := webhook.New(nil, nil, nil, "webhook_id", "webhook_token")

message, err := client.UpdateContent(context.ToDo(), "870741249114652722", "hello world!")
message, err := client.UpdateEmbeds(context.ToDo(), "870741249114652722", webhook.NewEmbedBuilder().
	SetDescription("hello world!").
	Build(),
)
message, err := client.UpdateMessage(context.ToDo(), "870741249114652722", webhook.NewMessageUpdateBuilder().
	SetContent("hello world!").
	Build(), 
)
```

### Delete Message

```go
client, err := webhook.New(nil, nil, nil, "webhook_id", "webhook_token")

err := client.DeleteMessage(context.ToDo(), "message_id")
```