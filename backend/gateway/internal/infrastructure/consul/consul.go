package consul

import (
	"gateway/internal/config"
	"gateway/internal/domain"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type Client struct {
	client      *consulapi.Client
	serviceName string
	srvList     *domain.ServerAvailableList

	logger *zap.Logger

	ticker *time.Ticker
}

func NewClient(cfg config.ConsulConfig, srvList *domain.ServerAvailableList, logger *zap.Logger) (*Client, error) {
	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfg.Addr

	client, err := consulapi.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:      client,
		serviceName: cfg.ServiceName,
		srvList:     srvList,
		logger:      logger,
		ticker:      time.NewTicker(time.Second * 10),
	}, nil
}

func (c *Client) HealthCheck() {
	for range c.ticker.C {
		health, _, err := c.client.Health().Service(c.serviceName, "", false, nil)
		if err != nil {
			c.logger.Error("cannot to observe available services via Consul", zap.Error(err))

			continue
		}

		var servers []string
		for _, item := range health {
			addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
			servers = append(servers, addr)
		}

		c.srvList.Update(servers)
	}
}

func (c *Client) Stop() {
	c.ticker.Stop()
}
