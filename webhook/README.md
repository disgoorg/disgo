# webhook

[Webhook](https://discord.com/developers/docs/resources/webhook) module of [disgo](https://github.com/DisgoOrg/disgo)

### Usage

Import the package into your project.

```go
import "github.com/DisgoOrg/disgo/webhook"
```

Create a new Webhook by `webhook_id` and `webhook_token`. (*This WebhookClient should be created once as it holds important state*)

```go
package main

import "github.com/DisgoOrg/disgo/webhook"

client := webhook.NewClient(discord.Snowflake(webhookID), webhookToken)
```

`webhook.NewClient` takes a vararg of type `webhook.ConfigOpt` as third argument which lets you pass additional optional parameter as a custom logger, rest client etc

### Optional Arguments

```go
package main

import "github.com/DisgoOrg/disgo/webhook"

client := webhook.NewClient(discord.Snowflake(webhookID), webhookToken,
	webhook.WithLogger(logrus.New()),
	webhook.WithDefaultAllowedMentions(discord.AllowedMentions{
        RepliedUser: false,
    }),
)
```

### Send Message

You can send a message as following

```go
package main

import (
    "github.com/DisgoOrg/disgo/core"
    "github.com/DisgoOrg/disgo/webhook"
)

client := webhook.NewClient(discord.Snowflake(webhookID), webhookToken)

message, err := client.CreateContent("hello world!")

message, err := client.CreateEmbeds(core.NewEmbedBuilder().
    SetDescription("hello world!").
    Build(),
)

message, err := client.CreateMessage(webhook.NewMessageCreateBuilder().
    SetContent("hello world!").
    Build(),
)
```

### Edit Message

Messages can also be edited

```go
package main

import (
    "github.com/DisgoOrg/disgo/core"
    "github.com/DisgoOrg/disgo/webhook"
)

client := webhook.NewClient(discord.Snowflake(webhookID), webhookToken)

message, err := client.UpdateContent("870741249114652722", "hello world!")

message, err := client.UpdateEmbeds("870741249114652722", core.NewEmbedBuilder().
    SetDescription("hello world!").
    Build(),
)

message, err := client.UpdateMessage("870741249114652722", webhook.NewMessageUpdateBuilder().
    SetContent("hello world!").
    Build(),
)
```

### Delete Message

or deleted

```go
package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/webhook"
)

client := webhook.NewClient(discord.Snowflake(webhookID), webhookToken)

err := client.DeleteMessage("message_id")
```