# voice

Voice provides a package to connect and send/receive voice to/from discord servers.
For Discords Docs on voice see [here](https://discord.com/developers/docs/topics/voice-connections).

DAVE(E2EE) library options:
 * https://github.com/disgoorg/godave CGO binding for https://github.com/discord/libdave
 * https://github.com/thomas-vilte/dave-go Pure Go implementation of DAVE(E2EE) (experimental)

## GoDave

### Installation

```bash
go get github.com/disgoorg/godave/golibdave
```

### Logging
Libdave uses a global logger which is set it `slog.LevelError` by default. You can change this by calling:

```go
libdave.SetDefaultLogLoggerLevel(slog.LevelInfo)
```

or set your own logger:

```go
libdave.SetDefaultLogLogger(yourLogger)
```

## Dave-Go

### Installation

```bash
go get github.com/thomas-vilte/dave-go
```

## Usage

To send audio you need to create a voice connection. When using the `bot.Client` package you can use `client.VoiceManager().CreateConn(guildID)`
```go
const (
    guildID = 12345
    channelID = 12345
)

client, err := disgo.New(token,
	bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
	bot.WithVoiceManagerConfigOpts(
		// for GoDave use this
		voice.WithDaveSessionCreateFunc(golibdave.NewSession),
		// for Dave-Go use this
		voice.WithDaveSessionCreateFunc(session.NewSession),
	),
)
// handle err

conn := client.VoiceManager().CreateConn(guildID)

err := conn.Open(context.TODO(), channelID, false, false)
// handle err

// set speaking flag
err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone)

// send opus frame
conn.UDP().Write(frame)

// close connection
conn.Close()
```

When using the voice package standalone you should create a voice manager. After this you can call `voice.Manager.CreateConn(guildID)`. After this you should send a `gateway.OpcodeVoiceStateUpdate` packet to the gateway.
```go
