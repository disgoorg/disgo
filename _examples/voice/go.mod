module github.com/disgoorg/disgo/_examples/application_commands/http

go 1.25.5

replace github.com/disgoorg/disgo => ../../

// FIXME: Remove before merging
replace github.com/disgoorg/godave => ../../../godave

require (
	github.com/disgoorg/disgo v0.19.0-rc14
	github.com/disgoorg/godave v0.0.1
	github.com/disgoorg/snowflake/v2 v2.0.3
)

require (
	github.com/disgoorg/json/v2 v2.0.0 // indirect
	github.com/disgoorg/omit v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)
