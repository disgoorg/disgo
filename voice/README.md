# voice

Voice provides a package to connect and send/receive voice to/from discord servers.
For Discords Docs on voice see [here](https://discord.com/developers/docs/topics/voice-connections).
Since DAVE(E2EE) will soon be required you also need https://github.com/disgoorg/godave

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
		voice.WithDaveSessionCreateFunc(golibdave.NewSession), 
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