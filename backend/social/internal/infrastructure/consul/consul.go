package consul

import (
	"fmt"
	"net"
	"social/internal/config"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

type Client struct {
	agent *consulapi.Agent
	cfg   config.ConsulConfig
}

func NewClient(cfg config.ConsulConfig) (*Client, error) {
	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfg.Addr

	client, err := consulapi.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		agent: client.Agent(),
		cfg:   cfg,
	}, nil
}

func (c *Client) Register() error {
	host, port, err := net.SplitHostPort(c.cfg.AgentAddr)
	if err != nil {
		return fmt.Errorf("fail to parse consul agent addr: %w", err)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("fail to parse consul agent port: %w", err)
	}

	if err = c.agent.ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      c.cfg.ServiceID,
		Name:    c.cfg.ServiceName,
		Port:    p,
		Address: host,
		Check: &consulapi.AgentServiceCheck{
			Interval: "5s",
			Timeout:  "3s",
			HTTP:     fmt.Sprintf("http://%s:%d/health-check", host, p),
		},
	}); err != nil {
		return fmt.Errorf("fail to sign up service via consul: %w", err)
	}

	return nil
}

func (c *Client) Deregister() error {
	if err := c.agent.ServiceDeregister(c.cfg.ServiceID); err != nil {
		return fmt.Errorf("fail to deregister service in consul: %w", err)
	}

	return nil
}
