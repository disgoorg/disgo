---
title: Slash Commands
prev: guide/disgo-bot/setting-up-your-bot/main-file
next: guide/disgo-bot/setting-up-your-bot/event-handling
weight: 3
---

[Slash commands](https://discord.com/developers/docs/interactions/application-commands) is the go-to way for users to interact with your bot on Discord.

They provide a lot of benefits over traditional message-based commands, such as:
- **Discoverability**: Users can easily find available commands by typing `/` in the chat.
- **Structured Input**: Commands can have predefined options and parameters, reducing user errors.
- **Permissions Handling**: Discord can automatically handle command permissions based on user roles (e.g. some commands are only available to moderators).
- **Autocomplete**: Commands can provide suggestions as users type, making it easier to use complex commands.

...and more!

## Creating your first Slash Command
Let's create the hello world equivalent of slash commands, aka "The Ping Pong Command".

The user will run the `/ping` command and the bot will respond with `Pong!`.

{{% steps %}}

### Define the Command
```go {filename="main.go"}
var (
  token   = os.Getenv("DISCORD_BOT_TOKEN")
  guildID = snowflake.GetEnv("DISCORD_GUILD_ID")
  commands = []discord.ApplicationCommandCreate{
    discord.SlashCommandCreate{
      Name:        "ping",
      Description: "Respond with Pong!",
    },
  }
)
```
A discord slash command is defined using the `discord.SlashCommandCreate` struct.<br>
It requires a `Name` and optionally a `Description`.<br>
`Name` is what the user will type after the `/` to invoke the command.

### Create the Command Handler
We defined our command, now we need to create a handler that will be called when a user uses our command.
```go {filename="main.go"}
func pingCommandHandler(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
  return e.CreateMessage(discord.MessageCreate{
    Content: "Pong!",
  })
}
```
This function listens for `ApplicationCommandInteractionCreate` events, which are triggered when a user uses an application command (Slash Command being one of them).<br>

### Register the Command
In order to use our command in Discord, we first need to register it with Discord's API.

You can register commands globally (available in all guilds your bot is in) or per-guild (only available in specific guilds).

We're going to take a look at how to register a command in a specific guild first, as it allows for faster testing and iteration. (*Also you don't want to deploy commands globally while you're still testing them!*)
```go {filename="main.go"}
//... inside main function, before creating the client
h := handler.New()
h.SlashCommand("/ping", pingCommandHandler)
if err := handler.SyncCommands(client, commands, []snowflake.ID{guildID}); err != nil {
  panic(err)
}

//... when creating the client
    bot.WithEventListeners(h),
```

This code snippet does a few things:
1. It creates a new handler using `handler.New()`.
2. It registers our `pingCommandHandler` function to handle the `/ping` command using `h.SlashCommand`.
3. It synchronizes our defined commands with Discord using `handler.SyncCommands`, specifying the guild ID where we want to register the command.
4. It adds the handler as an event listener to the bot client using `bot.WithEventListeners(h)`.

{{% /steps %}}

If you've been following along, your updated `main.go` file should now look like this:

```go {filename="main.go", linenos=table}
package main

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "time"

  // Import DisGo packages
  "github.com/disgoorg/disgo"
  "github.com/disgoorg/disgo/bot"
  "github.com/disgoorg/disgo/discord"
  "github.com/disgoorg/disgo/events"
  "github.com/disgoorg/disgo/gateway"
  "github.com/disgoorg/disgo/handler"
  "github.com/disgoorg/snowflake/v2"
)

var (
  token   = os.Getenv("DISCORD_BOT_TOKEN")
  guildID = snowflake.GetEnv("DISCORD_GUILD_ID")
  commands = []discord.ApplicationCommandCreate{
    discord.SlashCommandCreate{
      Name:        "ping",
      Description: "Respond with Pong!",
    },
  }
)

func pingCommandHandler(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
  return e.CreateMessage(discord.MessageCreate{
    Content: "Pong!",
  })
}

func main() {
  // Create a new client instance
  client, err := disgo.New(token,
    // Set gateway configuration options
    bot.WithGatewayConfigOpts(
      // Set enabled intents
      gateway.WithIntents(
        gateway.IntentGuilds,
        gateway.IntentGuildMessages,
        gateway.IntentDirectMessages,
      ),
    ),
    // Listen to the Ready event in order to know when the bot is connected
    bot.WithEventListenerFunc(func(e *events.Ready) {
      fmt.Println("Bot is connected as", e.User.Username)
    }),
  )
  if err != nil {
    panic(err)
  }

  h := handler.New()
  h.SlashCommand("/ping", pingCommandHandler)
  if err := handler.SyncCommands(client, commands, []snowflake.ID{guildID}); err != nil {
    panic(err)
  }

  // Ensure we close the client on exit
  defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client.Close(ctx)
  }()

  // Connect to the gateway
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  if err = client.OpenGateway(ctx); err != nil {
    panic(err)
  }

  // Wait here until CTRL+C or other term signal is received.
  fmt.Println("Bot is now running. Press CTRL+C to exit.")
  s := make(chan os.Signal, 1)
  signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
  <-s
  fmt.Println("Shutting down bot...")
}
```

Now try to run your bot and use the `/ping` command in your specified guild! You should see your bot respond with `Pong!`. 🏓

{{< callout type="info" emoji="💡">}}
  If you don't see the `/ping` command in your guild, try to refresh Discord (Ctrl + R / ⌘ + R) or wait a couple of minutes as sometimes it takes a bit of time for new commands to appear.
{{< /callout >}}

## Creating a Commands Directory
If you've been following this tutorial, our project directory should now look like this
{{< filetree/container >}}
  {{< filetree/folder name="awesome-discord-bot" >}}
    {{< filetree/file name="go.mod" >}}
    {{< filetree/file name="go.sum">}}
    {{< filetree/file name="main.go" >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

All our code is currently in the `main.go` file. While this is fine for small bots, as your bot grows in complexity, it's a good idea to start organizing your code into separate files and directories. This makes it easier to manage and maintain your codebase.

Make a new folder called `📂commands/` in the root of your project directory. This folder will hold all your command-related code.

Make a file in it called `commands.go` that will contain our slash command code.
```go {filename="commands/commands.go"}
package commands

import (
  "github.com/disgoorg/disgo/discord"
  "github.com/disgoorg/disgo/handler"
)

var Commands = []discord.ApplicationCommandCreate{
  discord.SlashCommandCreate{
    Name:        "ping",
    Description: "Respond with Pong!",
  },
}

func PingCommandHandler(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
  return e.CreateMessage(discord.MessageCreate{
    Content: "Pong!",
  })
}
```

Then, update your `main.go` file to use the new commands package.
```go {filename="main.go", hl_lines=[3, 7, 8]}
import(
  //... other imports
  "awesome-discord-bot/commands"
)

//... inside main function, after creating the client
h.SlashCommand("/ping", commands.PingCommandHandler)
if err := handler.SyncCommands(client, commands.Commands, []snowflake.ID{guildID}); err != nil {
  panic(err)
}
```
You can choose to add more commands in this file or create separate files for each command as your bot grows.
