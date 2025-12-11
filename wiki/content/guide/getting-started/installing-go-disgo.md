---
title: Installing Go & DisGo
prev: guide/getting-started/adding-bot
next: guide/disgo-bot/setting-up-your-bot/project-setup
weight: 3
---

I hope you've installed The Go programming language before following this tutorial (?)

If you have, you can skip to [here](#creating-your-disgo-project)

But if you haven't, I've got your back!

## Installing Go

Follow the installation steps listed [here](https://go.dev/doc/install)

## Creating your DisGo project

{{% steps %}}

### Go to your project directory

```bash
$ cd <your-bot-directory>
```

### Initialize Go module

```bash
$ go mod init awesome-discord-bot
```

### Install DisGo

```bash
$ go get github.com/disgoorg/disgo
```

You should now have an empty project with a single `go.mod` file that looks like this

```go {filename="go.mod"}
module awesome-discord-bot

go 1.25.4
```

{{% /steps %}}
