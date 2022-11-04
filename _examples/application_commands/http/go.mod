module github.com/disgoorg/disgo/_examples/application_commands/http

go 1.18

replace github.com/disgoorg/disgo => ../../../

require (
	github.com/disgoorg/disgo v0.11.5
	github.com/disgoorg/log v1.2.0
	github.com/disgoorg/snowflake/v2 v2.0.1
	github.com/oasisprotocol/curve25519-voi v0.0.0-20220317090546-adb2f9614b17
)

require (
	github.com/disgoorg/json v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e // indirect
	golang.org/x/exp v0.0.0-20220325121720-054d8573a5d8 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
)
