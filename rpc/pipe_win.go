//go:build windows

package rpc

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

func GetDiscordIPCPath(i int) string {
	return fmt.Sprintf("\\\\?\\pipe\\discord-ipc-%d", i)
}

func openPipe(path string) (net.Conn, error) {
	return winio.DialPipe(path, nil)
}
