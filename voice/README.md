# voice

Voice provides a package to connect and send/receive voice to/from discord servers.
For Discords Docs on voice see [here](https://discord.com/developers/docs/topics/voice-connections).

## Usage

To send audio you need to create a voice connection. When using the `bot.Client` package you can use `bot.OpenVoice()`
```go
const (
    guildID = 12345
    channelID = 12345
)

var client bot.Client

conn, err := client.OpenVoice(context.TODO(), guildID, channelID)
```

When using the voice package standalone you should create a voice manager. After this you can call `voice.Manager.CreateConn()`. After this you should send a `gateway.OpcodeVoiceStateUpdate` packet to the gateway.
```go