package config

import (
	"fmt"
	"os"
	"strconv"
)

type ConfigNATS struct {
	Host    string
	Port    int
	Client  string
	Cluster string
	Channel string
	Queue   string
	Durable string
}

func (c *ConfigNATS) Getenv() error {
	var err error
	c.Port, err = strconv.Atoi(os.Getenv("NATS_PORT"))
	if err != nil {
		return err
	}
	c.Host = os.Getenv("NATS_HOST")
	c.Client = os.Getenv("NATS_CLIENT")
	c.Cluster = os.Getenv("NATS_CLUSTER")
	c.Channel = os.Getenv("NATS_CHANNEL")
	c.Queue = os.Getenv("NATS_QUEUE")
	c.Durable = os.Getenv("NATS_DURABLE")
	return nil
}

func (c *ConfigNATS) GetConnStr() string {
	return fmt.Sprintf("nats://%s:%d", c.Host, c.Port)
}
