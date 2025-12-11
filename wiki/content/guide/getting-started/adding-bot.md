---
title: Adding your Bot
prev: guide/getting-started/app-setup
next: guide/getting-started/installing-go-disgo
weight: 2
params:
  images:
    - "/images/getting-started/adding-your-bot.png"
---

Remember our bot? It's now alive, but it's sitting in Discord jail and not in any server like it's supposed to, cause we haven't invited it to any server yet.
So let's do that now!

## Bot invite links

Head over to the "OAuth2" section in the sidebar on the left
![](/images/getting-started/OAuth2.png)

Scroll down to the "OAuth2 URL Generator"
![](/images/getting-started/OAuth2_url_generator.png)

Here you'll find checkboxes where you can specify what permissions you want your bot to have when you invite it to a server and the scope that it should operate on. Our most important one is `bot`, `applications.commands`, & `Administrator` for this tutorial as that is what will allow us to invite our App as a Discord Bot and give it permission to be able to run commands.

{{< callout type="info" >}}
It's usually not recommended to give a Bot `Administrator` permission unless you actually need it, so it's recommended to fine tune the scope & the permissions your bot needs beyond this tutorial
{{< /callout >}}

Finally, select either "Guild Install" (if your bot should be run within a guild) or "User Install" (if your bot can be used by a user personally in a server or in DMs)

Copy the generated URL, select which server you want to install your bot to and **Congratulations!**<br>
You're finally ready to code an awesome discord bot using DisGo
