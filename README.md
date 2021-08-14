[![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/disgo.svg)](https://pkg.go.dev/github.com/DisgoOrg/disgo)
[![Go Report](https://goreportcard.com/badge/github.com/DisgoOrg/disgo)](https://goreportcard.com/report/github.com/DisgoOrg/disgo)
[![Go Version](https://img.shields.io/github/go-mod/go-version/DisgoOrg/disgo)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgo/blob/master/LICENSE)
[![Disgo Version](https://img.shields.io/github/v/tag/DisgoOrg/disgo?label=release)](https://github.com/DisgoOrg/disgo/releases/latest)
[![Disgo Discord](https://img.shields.io/discord/817327181659111454?color=%231081c1&label=discord)](https://discord.gg/TewhTfDpvW)

<img align="right" src="/.github/discord_gopher.png" width=192 alt="discord gopher">

# disgo

disgo is a [Discord](https://discord.com/developers/docs/resources/webhook) API wrapper written
in [Go](https://golang.org/) aimed to be consistent, modular, customizable and easy to use

## Summary

1. [Features](#features)
2. [Missing Features](#missing-features)
3. [Getting Started](#getting-started)
4. [Documentation](#documentation)
5. [Examples](#examples)
6. [Related Projects](#related-projects)
7. [Troubleshooting](#troubleshooting)
8. [Contributing](#contributing)
9. [License](#license)

### Features

* Interactions over [HTTP](https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction)
  support
* Full [Slash Commands](https://discord.com/developers/docs/interactions/slash-commands) support
* Full [Message Components](https://discord.com/developers/docs/interactions/message-components) support
* Full Rest API coverage
* [Gateway](https://discord.com/developers/docs/topics/gateway) support
* Full [Stage Instance](https://discord.com/developers/docs/resources/stage-instance) support
* [Guild Template](https://discord.com/developers/docs/resources/guild-template) support
* [RateLimit](https://discord.com/developers/docs/topics/rate-limits) handling
* [Webhook](https://discord.com/developers/docs/resources/webhook) support

### Missing Features

* [Voice](https://discord.com/developers/docs/resources/voice) support
* [Sharding](https://discord.com/developers/docs/topics/gateway#sharding) support
* [Threads](https://discord.com/developers/docs/topics/threads) support
* [OAuth2](https://discord.com/developers/docs/topics/oauth2) support

## Getting Started

### Installing

```sh
go get github.com/DisgoOrg/disgo
```

### Building a Disgo Instance

```go
disgo, err := core.NewBuilder("token").
    // set which gateway intents we should use
    SetGatewayConfig(gateway.Config{
            GatewayIntents: gateway.IntentGuilds | gateway.IntentGuildMessages,
    }).
    SetHTTPServerConfig(httpserver.Config{
        URL:       "/interactions/callback",
        Port:      ":443",
        PublicKey: "your public key from the developer dashboard",
    }).
    // add our event listeners
    AddEventListeners(&events.ListenerAdapter{
        OnGuildMessageCreate: guildMessageListener,
    }).
    // build the disgo instance. This might return an error!
    Build()

// connect to the gateway
err := disgo.Connect()

// optionally start the http server for interactions
disgo.Start() 
```

### Ping Pong Example

```go
package main

import (
    "os"
    "os/signal"
    "syscall"

    "github.com/DisgoOrg/disgo/core"
    "github.com/DisgoOrg/disgo/core/events"
    "github.com/DisgoOrg/disgo/gateway"
    "github.com/DisgoOrg/log"
)

func main() {
    // create a new builder
    disgo, err := core.NewBuilder("token").
    	// set which gateway intents we should use
        SetGatewayConfig(gateway.Config{
                GatewayIntents: gateway.IntentGuilds | gateway.IntentGuildMessages,
        }).
    	// add our event listeners
        AddEventListeners(&events.ListenerAdapter{
            OnGuildMessageCreate: guildMessageListener,
        }).
    	// build the disgo instance. This might return an error!
        Build()
    if err != nil {
        log.Fatal("error while building disgo: ", err)
    }

    // clean exit disgo 
    defer disgo.Close()

    // connect to the gateway to receive events from discord
    if err = disgo.Connect(); err != nil {
        log.Fatal("failed to connect to gateway: ", err)
    }
    
    // block until we receive a stop signal
    s := make(chan os.Signal, 1)
    signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-s
}

// define event listener function to get message create events
func guildMessageListener(event *events.GuildMessageCreateEvent) {
    message := event.Message
    // check if message author is bot or content is nil
    if message.Author.IsBot || message.Content == nil {
        return
    }

    // check if message content is "ping"
    if *message.Content == "ping" {
    	// reply to the message with pong
        if _, err := message.Reply(core.NewMessageCreateBuilder().
            SetContent("pong").
            Build(),
        ); err != nil {
            log.Error("failed to reply to ping: ", err)
        }
    }
}
```

### Logging

disgo uses our own small [logging lib](https://github.com/DisgoOrg/log) which provides
an [interface](https://github.com/DisgoOrg/log/blob/master/logger.go) you can implement. This lib also comes with a
default logger which is interchangeable and based on the standard log package. You can read more about
it [here](https://github.com/DisgoOrg/log)

## Documentation

Documentation is unfinished and can be found under

* [![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/disgo.svg)](https://pkg.go.dev/github.com/DisgoOrg/disgo)
* [![Discord Documentation](https://img.shields.io/badge/Discord%20Documentation-blue.svg)](https://discord.com/developers/docs)

Wiki is currently under construction

## Examples

You can find examples under [_examples](https://github.com/DisgoOrg/disgo/tree/master/_examples)

or in these projects:

* [disgo-butler](https://github.com/DisgoOrg/disgo-butler)
* [BansBot](https://github.com/Skye-31/BansBot)
* [Reddit-Discord-Bot](https://github.com/TopiSenpai/Reddit-Discord-Bot)
* [Kitsune-Bot](https://github.com/TopiSenpai/Kitsune-Bot)
* [Uno-Bot](https://github.com/TopiSenpai/Uno-Bot)

## Related Projects

### [Lavalink](https://github.com/freyacodes/Lavalink)

Standalone audio sending node based on Lavaplayer and JDA-Audio. Allows for sending audio without it ever reaching any
of your shards.

Being used in production by FredBoat, Dyno, LewdBot, and more.

### [disgolink](https://github.com/DisgoOrg/disgolink)

[Lavalink Client](https://github.com/freyacodes/Lavalink) which can be used to communicate with LavaLink to play/search
tracks

### [disgofy](https://github.com/DisgoOrg/disgofy)

[disgolink](https://github.com/DisgoOrg/disgolink) Spotify integration. disgofy resolved Spotify urls to
tracks/albums/playlists and lazy searches for them on YouTube

### [dislog](https://github.com/DisgoOrg/dislog)

Discord webhook logger integration for [logrus](https://github.com/sirupsen/logrus)

### [disgommand](https://github.com/DisgoOrg/disgommand)

Command framework for disgo in [gorilla/mux](https://github.com/gorilla/mux) style

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/TewhTfDpvW)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/TewhTfDpvW) or
create an issue to discuss your intentions and ideas.

## License

Distributed under
the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgo/blob/master/LICENSE)
. See LICENSE for more information.

