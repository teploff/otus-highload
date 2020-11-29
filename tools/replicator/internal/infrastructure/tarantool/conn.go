package tarantool

import (
	"github.com/tarantool/go-tarantool"
	"replicator/internal/config"
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

func (c *Conn) Insert(values ...interface{}) error {
	row := make([]interface{}, 0, len(values))
	for _, v := range values {
		row = append(row, v)
	}

	_, err := c.conn.Insert(c.space, row)

	return err
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
