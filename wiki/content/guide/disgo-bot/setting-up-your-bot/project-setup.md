---
title: Project Setup
prev: guide/getting-started/installing-go-disgo
next: guide/disgo-bot/setting-up-your-bot/main-file
weight: 1
---

## Configuration Files

Remember your [Token?](/guide/getting-started/app-setup/#the-bot-token)
We're going to need it now.

Recall that the token is **Very important**, and we need to protect it from being leaked as best as we possibly can

One of the ways of doing that is to have a configuration file for our bot that will have our token, and **DO NOT COMMIT THIS FILE INTO A PUBLIC REPOSITORY**

You can set up a config file in `config.json`, `config.yaml`, / `config.toml` file, you can also use a `.env` file, or just supply it directly from the terminal when running your command, the choice is up to you

An example of a config file looks like this
{{< tabs items="JSON, .env" >}}
  {{< tab >}}

  ```json {filename="config.json"}
  {
    "DISCORD_TOKEN": "insert-your-token-here"
  }
  ```

  {{< /tab >}}

  {{< tab >}}

  ```env {filename=".env"}
  DISCORD_TOKEN=insert-your-token-here
  ```

  {{< /tab >}}
{{< /tabs >}}

Whichever way you choose to do it, just make sure that it's in your `.gitignore` file so that it doesn't get committed
