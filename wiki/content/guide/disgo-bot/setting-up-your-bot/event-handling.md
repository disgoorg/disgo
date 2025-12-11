---
title: Event Handling
prev: guide/disgo-bot/setting-up-your-bot/slash-commands
weight: 4
---

So far we've only used commands to interact with our bot, but our bot can do a lot more than just respond to commands!

## What are Events?
Events are actions that happen in Discord, such as a user sending a message, a user joining a guild, a reaction being added to a message, etc.

Actually, we've used 2 of them already!<br>
- When a user uses a Slash Command, it triggers an `ApplicationCommandInteractionCreate` event, which our bot listens for and responds to.<br>
- When we start our bot, it also triggers a `Ready` event, which we can use to perform actions when the bot is fully connected and ready.

Your bot can listen for these events and respond to them accordingly.

## Setting up Event Handlers
To handle events in DisGo, we can use the `events` package from DisGo, then register it to our bot using `bot.WithEventListeners` / `bot.WithEventListenerFunc` when creating the bot client.

Here's an example of replying to a user if they mention our bot (Note that you'll need `Message Content` intent enabled in your bot settings for this to work)
```go {filename="main.go"}
func onMessageCreate(e *events.MessageCreate) {
  if e.Message.Author.Bot {
    return
  }

  if strings.Contains(e.Message.Content, e.Client().ApplicationID().String()) {
    e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
      Content: "You called? 📞",
    })
  }
}

//... when creating the client
    bot.WithEventListenerFunc(onMessageCreate),
```
![](/images/disgo-bot/setting-up-your-bot/message-create-example.png)

You can see a list of all available events in the [DisGo Events Package](https://pkg.go.dev/github.com/disgoorg/disgo/events).
## Further Reading
That's mostly it for our basic bot setup! From here you can explore more advanced topics such as:

<div class="hx:mt-6 hx:mb-6">
  {{< cards >}}
    {{< card title="Message Components" subtitle="Display interactive buttons, dropdowns, and more" icon="template" >}}
    {{< card title="Slash Commands Config" subtitle="Create commands that takes options, autocomplete, and more" icon="adjustments" >}}
    {{< card title="Command response methods" subtitle="Learn about different ways to respond to commands" icon="at-symbol" >}}
    {{< card title="App Sharding" subtitle="Scale your bot with multiple shards" icon="server" >}}
  {{< /cards >}}
</div>
