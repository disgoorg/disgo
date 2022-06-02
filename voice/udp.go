package voice

import "net"

func NewUDPConn(conn net.Conn, opts ...UDPConnConfigOpt) *UDPConn {
	config := DefaultUDPConnConfig()
	config.Apply(opts)

	return &UDPConn{
		Conn:   conn,
		config: *config,
	}
}

type UDPConn struct {
	Conn net.Conn
}

func (c *UDPConn) Close() {
	_ = c.Conn.Close()
}
