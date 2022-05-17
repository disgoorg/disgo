[![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/disgo.svg)](https://pkg.go.dev/github.com/disgoorg/disgo)
[![Go Report](https://goreportcard.com/badge/github.com/disgoorg/disgo)](https://goreportcard.com/report/github.com/disgoorg/disgo)
[![Go Version](https://img.shields.io/github/go-mod/go-version/disgoorg/disgo)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgo/blob/master/LICENSE)
[![Disgo Version](https://img.shields.io/github/v/tag/disgoorg/disgo?label=release)](https://github.com/disgoorg/disgo/releases/latest)
[![Disgo Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/TewhTfDpvW)

<img align="right" src="/.github/discord_gopher.png" width=192 alt="discord gopher">

# DisGo

DisGo is a [Discord](https://discord.com) API wrapper written in [GoLang](https://golang.org/) aimed to be consistent, modular, customizable and higher level than other Discord API wrappers.

## Summary

1. [Stability](#stability)
2. [Features](#features)
3. [Missing Features](#missing-features)
4. [Getting Started](#getting-started)
5. [Documentation](#documentation)
6. [Examples](#examples)
7. [Other interesting Projects to look at](#other-interesting-projects-to-look-at)
8. [Other GoLang Discord Libraries](#other-golang-discord-libraries)
9. [Troubleshooting](#troubleshooting)
10. [Contributing](#contributing)
11. [License](#license)

### Stability
The public API of DisGo is mostly stable at this point in time. Smaller breaking changes can happen before the v1 is released. 

After v1 is released breaking changes may only happen if the Discord API requires them. They tend to break their released API versions now and then. In general for every new Discord APi version the major version of DisGo should be increased and with that breaking changes between non major versions should be helt to a minimum. 

### Features

* Full Rest API coverage
* [Gateway](https://discord.com/developers/docs/topics/gateway)
* [Sharding](https://discord.com/developers/docs/topics/gateway#sharding)
* [HTTP Interactions](https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction)
* [Application Commands](https://discord.com/developers/docs/interactions/application-commands)
* [Message Components](https://discord.com/developers/docs/interactions/message-components)
* [Modals](https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-modal)
* [Stage Instance](https://discord.com/developers/docs/resources/stage-instance)
* [Guild Template](https://discord.com/developers/docs/resources/guild-template)
* [Sticker](https://discord.com/developers/docs/resources/sticker)
* [RateLimit](https://discord.com/developers/docs/topics/rate-limits)
* [Webhook](https://discord.com/developers/docs/resources/webhook)
* [OAuth2](https://discord.com/developers/docs/topics/oauth2)
* [Threads](https://discord.com/developers/docs/topics/threads)
* [Guild Scheduled Event](https://discord.com/developers/docs/resources/guild-scheduled-event)

### Missing Features

* [Voice](https://discord.com/developers/docs/topics/voice-connections)
* [RPC](https://discord.com/developers/docs/topics/rpc)

## Getting Started

### Installing

```sh
go get github.com/disgoorg/disgo
```

### Building a DisGo Instance

Build a bot client to interact with the discord api
```go
package main

import (
    "github.com/disgoorg/disgo"
    "github.com/disgoorg/disgo/bot"
    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/events"
    "github.com/disgoorg/disgo/gateway"
)

func main() {
    client, err := disgo.New("token",
        // set gateway options
        bot.WithGatewayConfigOpts(
            // set enabled intents
            gateway.WithGatewayIntents(
                discord.GatewayIntentGuilds,
                discord.GatewayIntentGuildMessages,
                discord.GatewayIntentDirectMessages,
            ),
        ),
        // add event listeners
        bot.WithEventListenerFunc(func(e *events.MessageCreateEvent) {
            // event code here
        }),
    )
}
```

### Full Ping Pong Example

This example can also be found [here]https://github.com/disgoorg/disgo/blob/development/_examples/ping_pong/example.go()

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/disgoorg/disgo"
    "github.com/disgoorg/disgo/bot"
    "github.com/disgoorg/disgo/cache"
    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/events"
    "github.com/disgoorg/disgo/gateway"
    "github.com/disgoorg/log"
)

func main() {
    client, err := disgo.New(os.Getenv("token"),
        bot.WithGatewayConfigOpts(
            gateway.WithGatewayIntents(
                discord.GatewayIntentGuildMessages,
                discord.GatewayIntentMessageContent,
            ),
        ),
        bot.WithEventListenerFunc(onMessageCreate),
    )
    if err != nil {
        log.Fatal("error while building disgo: ", err)
    }

    defer client.Close(context.TODO())

    if err = client.ConnectGateway(context.TODO()); err != nil {
        log.Fatal("errors while connecting to gateway: ", err)
    }

    log.Info("example is now running. Press CTRL-C to exit.")
    s := make(chan os.Signal, 1)
    signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
    var message string
    if event.Message.Content == "ping" {
        message = "pong"
    } else if event.Message.Content == "pong" {
        message = "ping"
    }
    if message != "" {
        _, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
    }
}
```

### Logging

DisGo uses our own small [logging interface](https://github.com/disgoorg/log) which you can use with most other logging libaries. This lib also comes with a default logger which is based on the standard log package.

## Documentation

Documentation is wip and can be found under

* [![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/disgo.svg)](https://pkg.go.dev/github.com/disgoorg/disgo)
* [![Discord Documentation](https://img.shields.io/badge/Discord%20Documentation-blue.svg)](https://discord.com/developers/docs)

GitHub Wiki is currently under construction. We appreciate help here.

## Examples

You can find examples [here](https://github.com/disgoorg/disgo/tree/master/_examples)

There is also a bot template with commands & db [here](https://github.com/disgoorg/bot-template)

or in these projects:

* [disgo-butler](https://github.com/disgoorg/disgo-butler)
* [BansBot](https://github.com/Skye-31/BansBot)
* [Reddit-Discord-Bot](https://github.com/TopiSenpai/Reddit-Discord-Bot)
* [Kitsune-Bot](https://github.com/TopiSenpai/Kitsune-Bot)
* [KittyBot](https://github.com/KittyBot-Org/KittyBotGo)

## Other interesting Projects to look at

### [Lavalink](https://github.com/freyacodes/Lavalink)

Is a standalone audio sending node based on [Lavaplayer](https://github.com/sedmelluq/lavaplayer) and JDA-Audio. Which allows for sending audio without it ever reaching any of your shards. Lavalink can be used in combinatio with [DisGolink](https://github.com/disgoorg/disgolink) for music Bots

Being used in production by FredBoat, Dyno, LewdBot, and more.

### [DisGolink](https://github.com/disgoorg/disgolink)

Is a [Lavalink-Client](https://github.com/freyacodes/Lavalink) which can be used to communicate with LavaLink to play/search tracks

### [DisLog](https://github.com/disgoorg/dislog)

Is a Discord webhook logger hook for [logrus](https://github.com/sirupsen/logrus)

## Other GoLang Discord Libraries

* [discordgo](https://github.com/bwmarrin/discordgo)
* [disgord](https://github.com/andersfylling/disgord)
* [arikawa](https://github.com/diamondburned/arikawa)
* [corde](https://github.com/Karitham/corde)

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/TewhTfDpvW)

## Contributing

Contributions are welcomed but for bigger changes we recommend first reaching out via [Discord](https://discord.gg/TewhTfDpvW) or create an issue to discuss your problems, intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgo/blob/master/LICENSE)
. See LICENSE for more information.


