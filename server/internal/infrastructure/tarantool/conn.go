package tarantool

import (
	"github.com/tarantool/go-tarantool"
	"social-network/internal/config"
	"strconv"
	"time"
)

type Conn struct {
	conn  *tarantool.Connection
	space string
}

func NewConn(config config.TarantoolConfig) (*Conn, error) {
	c, err := tarantool.Connect(config.Host+":"+strconv.Itoa(config.Port), tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          config.User,
		Pass:          config.Password,
		Logger:        nil,
	})
	if err != nil {
		return nil, err
	}

	_, err = c.Ping()
	if err != nil {
		return nil, err
	}

	return &Conn{conn: c, space: config.Space}, nil
}

func (c *Conn) CallFunc(functionName string, args interface{}) ([]interface{}, error) {
	resp, err := c.conn.Call(functionName, args)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
