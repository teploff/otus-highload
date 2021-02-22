package zabbix

import (
	"context"
	"messenger/internal/config"
	"time"

	zbx "github.com/blacked/go-zabbix"
	"go.uber.org/zap"
)

type Client struct {
	sender   *zbx.Sender
	hostName string

	logger *zap.Logger
}

func NewClient(cfg config.ZabbixConfig, logger *zap.Logger) *Client {
	return &Client{
		sender:   zbx.NewSender(cfg.ServerHost, cfg.Port),
		hostName: cfg.HostName,
		logger:   logger,
	}
}

func (c *Client) Publish(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pkg := c.collectMetrics()

			if _, err := c.sender.Send(pkg); err != nil {
				c.logger.Error("send pkg to zabbix", zap.Error(err))
			}
		}
	}

}

func (c *Client) collectMetrics() *zbx.Packet {
	metrics := make([]*zbx.Metric, 0, 1)

	metrics = append(metrics, zbx.NewMetric(c.hostName, "messenger-cpu", "1.22", time.Now().Unix()))
	metrics = append(metrics, zbx.NewMetric(c.hostName, "status", "OK", time.Now().Unix()))

	return zbx.NewPacket(metrics)
}
