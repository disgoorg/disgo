//go:build !windows

package rpc

import (
	"fmt"
	"net"
	"os"
)

var paths = []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

func getDiscordIPCPath(i int) string {
	tmpPath := "/tmp"
	for _, path := range paths {
		if v := os.Getenv(path); v != "" {
			tmpPath = v
			break
		}
	}

	if tmpPath[len(tmpPath)-1] != '/' {
		tmpPath += "/"
	}
	// TODO: support tmpPath/snap.discord/discord-ipc-n
	// TODO: support tmpPath/app/com.discordapp.Discord/discord-ipc-n
	return fmt.Sprintf("%sdiscord-ipc-%d", tmpPath, i)
}

func openPipe(path string) (net.Conn, error) {
	return net.Dial("unix", path)
}
