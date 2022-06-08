module github.com/disgoorg/disgo/_examples/audio

go 1.18

replace (
	github.com/disgoorg/disgo => ../../
	github.com/disgoorg/disgo/audio => ../../audio
)

require (
	github.com/disgoorg/disgo v0.12.3
	github.com/disgoorg/disgo/audio v0.12.1
	github.com/disgoorg/log v1.2.0
)

require (
	github.com/disgoorg/snowflake/v2 v2.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/exp v0.0.0-20220325121720-054d8573a5d8 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
)
