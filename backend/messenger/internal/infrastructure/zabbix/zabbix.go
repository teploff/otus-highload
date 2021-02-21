package zabbix

import (
	"context"
	"messenger/internal/config"
	"time"

	zbx "github.com/blacked/go-zabbix"
	"go.uber.org/zap"
)

type Client struct {
	sender *zbx.Sender

	logger *zap.Logger
}

func NewClient(cfg config.ZabbixConfig, logger *zap.Logger) *Client {
	return &Client{
		sender: zbx.NewSender(cfg.Host, cfg.Port),
		logger: logger,
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

			res, err := c.sender.Send(pkg)
			if err != nil {
				c.logger.Error("send pkg to zabbix", zap.Error(err))
			}

			c.logger.Info(string(res))
		}
	}

}

func (c *Client) collectMetrics() *zbx.Packet {
	metrics := make([]*zbx.Metric, 0, 1)

	metrics = append(metrics, zbx.NewMetric("messenger", "messenger-cpu", "1.22", time.Now().Unix()))
	metrics = append(metrics, zbx.NewMetric("messenger", "status", "OK", time.Now().Unix()))

	return zbx.NewPacket(metrics)
}
