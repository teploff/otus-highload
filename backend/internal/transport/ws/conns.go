package ws

import (
	"net"
	"sync"
)

type Conns struct {
	conns map[string][]net.Conn

	m sync.Mutex
}

func NewWSConnects() *Conns {
	return &Conns{
		conns: make(map[string][]net.Conn),
		m:     sync.Mutex{},
	}
}

func (c *Conns) Add(userID string, conn net.Conn) {
	c.m.Lock()
	defer c.m.Unlock()

	c.conns[userID] = append(c.conns[userID], conn)
}

func (c *Conns) Remove(userID string, conn net.Conn) {
	c.m.Lock()
	defer c.m.Unlock()

	for index, connection := range c.conns[userID] {
		if connection == conn {
			c.conns[userID][index] = c.conns[userID][len(c.conns[userID])-1]
			c.conns[userID][len(c.conns[userID])-1] = nil
			c.conns[userID] = c.conns[userID][:len(c.conns[userID])-1]

			return
		}
	}
}

func (c *Conns) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	for _, conns := range c.conns {
		for _, conn := range conns {
			conn.Close()
		}
	}
}
