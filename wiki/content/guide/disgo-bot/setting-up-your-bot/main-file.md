---
title: Main File
prev: guide/disgo-bot/setting-up-your-bot/project-setup
next: guide/disgo-bot/setting-up-your-bot/slash-commands
weight: 2
---

So far, what we've had as a project looks like this

{{< filetree/container >}}
  {{< filetree/folder name="awesome-discord-bot" >}}
    {{< filetree/file name="go.mod" >}}
    {{< filetree/file name="go.sum">}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

Not so much going on huh?

## main.go

Add a new file called `main.go`, this is our main entrypoint for our bot.

Here's a code snippet you can use to get started:

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
  "github.com/disgoorg/snowflake/v2"
)

var (
  token   = os.Getenv("DISCORD_BOT_TOKEN")
  guildID = snowflake.GetEnv("DISCORD_GUILD_ID")
)

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

That's a lot of code, but don't worry if you don't understand everything right now.
We'll break it down

First we create a new DisGo client using our bot token
```go {filename="main.go", linenos=table, linenostart=16, hl_lines=[3]}
func main() {
  // Create a new client instance
  client, err := disgo.New(token,
  //...

```
This creates a new DisGo client that will provide all the functionality we need to interact with the Discord API.

Next we set some options for the gateway connection, including which intents we want to enable.
```go {filename="main.go", linenos=table, linenostart=17}
    // Set gateway configuration options
    bot.WithGatewayConfigOpts(
      // Set enabled intents (https://discord.com/developers/docs/events/gateway#gateway-intents)
      gateway.WithIntents(
        gateway.IntentGuilds,
        gateway.IntentGuildMessages,
        gateway.IntentDirectMessages,
      ),
    ),
    // Listen to the Ready event in order to know when the bot is connected
    bot.WithEventListenerFunc(func(e *events.Ready) {
      fmt.Print("Bot is connected as", e.User.Username)
    }),
```
Here we specify that we want to receive events related to guilds (servers), guild messages, and direct messages.

We've also set up an event listener for the `Ready` event, which is triggered when the bot successfully connects to Discord.

After creating the client, we set up a deferred function to ensure that the client is properly closed when the application exits:
```go {filename="main.go", linenos=table, linenostart=38}
// Ensure we close the client on exit
  defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client.Close(ctx)
  }()
```

We then try to connect to the gateway
```go {filename="main.go", linenos=table, linenostart=44, hl_lines=[4]}
  // Connect to the gateway
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  if err = client.OpenGateway(ctx); err != nil {
    panic(err)
  }
```

And if everything works correctly, our bot is now ready to run and connect to Discord!

## Running your application

If you've been following along, your project structure should now look like this

{{< filetree/container >}}
  {{< filetree/folder name="awesome-discord-bot" >}}
    {{< filetree/file name="go.mod" >}}
    {{< filetree/file name="go.sum">}}
    {{< filetree/file name="main.go" >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

Head to the terminal and run the following command inside your project directory to run your bot

```bash
go run main.go
```

If everything is set up correctly, you should see a message indicating that your bot is connected with the username you've set for it.
```
Bot is Connected as YourBotName
```

Congratulations! You've successfully set up the main file for your DisGo bot and connected it to Discord.

{{< callout type="error" emoji="⚠️" >}}
  **Remember**, never share your bot token with anyone or expose it in public repositories. Keep it secure!
{{< /callout >}}

Now we have a bot, but it doesn't do much yet. In the next section, we'll start adding some functionality to our bot!
