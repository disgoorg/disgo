module github.com/disgoorg/disgo/_examples/application_commands/http

go 1.18

replace github.com/disgoorg/disgo => ../../../

require (
	github.com/disgoorg/disgo v0.7.4
	github.com/disgoorg/log v1.2.0
	github.com/disgoorg/snowflake v1.1.0
	github.com/oasisprotocol/curve25519-voi v0.0.0-20220317090546-adb2f9614b17
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)
