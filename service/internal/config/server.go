package config

import (
	"fmt"
	"os"
	"strconv"
)

type ConfigServer struct {
	Host string
	Port int
}

func (c *ConfigServer) Getenv() error {
	var err error
	c.Port, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return err
	}
	c.Host = os.Getenv("SERVER_HOST")
	return nil
}

func (c *ConfigServer) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
