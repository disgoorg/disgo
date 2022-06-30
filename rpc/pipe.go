package rpc

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

var unixPaths = []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

func DialPipe() (*Conn, error) {
	conn, err := openPipe(GetDiscordIPCPath(0))
	if conn != nil {
		return conn, nil
	}
	return nil, err
}

func GetDiscordIPCPath(i int) string {
	if strings.HasPrefix(runtime.GOOS, "windows") {
		return fmt.Sprintf("\\\\?\\pipe\\discord-ipc-%d", i)
	}

	tmpPath := "/tmp"
	for _, path := range unixPaths {
		if v := os.Getenv(path); v != "" {
			tmpPath = v
			break
		}
	}
	return fmt.Sprintf("%sdiscord-ipc-%d", tmpPath, i)
}

func openPipe(path string) (*Conn, error) {
	var (
		conn net.Conn
		err  error
	)

	if strings.HasPrefix(runtime.GOOS, "windows") {
		return nil, nil
	} else {
		conn, err = net.Dial("unix", path)
	}

	if err != nil {
		return nil, err
	}
	return newConn(conn), nil
}
