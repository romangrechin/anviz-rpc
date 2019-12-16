package device

import (
	"github.com/romangrechin/anviz-rpc/anviz/errors"
	"net"
	"time"
)

const (
	DISCONNECTED = iota
	CONNECTED
	BUSY
)

var (
	emptyTime = time.Time{}
)

type connection struct {
	id               uint32
	address          string
	conn             net.Conn
	isBusy           bool
	maxReconnect     int32
	currentReconnect int32
	reconnectTimeout int32
	status           uint8
	done             chan struct{}
	connTimeout      time.Duration
	readWriteTimeout time.Duration
}

func (c *connection) IsBusy() bool {
	return c.isBusy
}

func (c *connection) setIsBusy(val bool) {
	c.isBusy = val
}

func (c *connection) dial() (err error) {
	c.setIsBusy(true)
	defer c.setIsBusy(false)
	c.conn, err = net.DialTimeout("tcp", c.address, c.connTimeout)
	if err != nil {
		err = errors.ErrCouldNotConnect
		return
	}
	c.status = CONNECTED
	return
}

func (c *connection) send(cmd uint8, data []byte) ([]byte, error) {
	if c.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}

	if c.status != CONNECTED {
		return nil, errors.ErrConnectionClosed
	}

	c.setIsBusy(true)
	defer c.setIsBusy(false)

	buf := marshal(cmd, c.id, data)
	if c.readWriteTimeout > 0 {
		c.conn.SetDeadline(time.Now().Add(c.readWriteTimeout))
	}
	_, err := c.conn.Write(buf)
	if err != nil {
		c.Close()
		return nil, errors.ErrConnectionWrite
	}

	resBuf := make([]byte, 512)
	if c.readWriteTimeout > 0 {
		c.conn.SetDeadline(time.Now().Add(c.readWriteTimeout))
	}
	n, err := c.conn.Read(resBuf)
	if err != nil {
		c.Close()
		return nil, errors.ErrConnectionRead
	}
	c.conn.SetDeadline(emptyTime)
	id, response, err := unmarshal(cmd, resBuf[:n])
	c.id = id
	return response, err
}

func (c *connection) Close() {
	if c.conn != nil {
		c.status = DISCONNECTED
		_ = c.conn.Close()
	}
}

func newConnection(address string) (*connection, error) {
	conn := &connection{
		address: address,
	}

	err := conn.dial()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
