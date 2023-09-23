package config

import (
	"fmt"
	"os"
	"strconv"
)

type ConfigDB struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

func (c *ConfigDB) Getenv() error {
	var err error
	c.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return err
	}
	c.Host = os.Getenv("DB_HOST")
	c.User = os.Getenv("DB_USER")
	c.Pass = os.Getenv("DB_PASS")
	c.Name = os.Getenv("DB_NAME")
	return nil
}

func (c *ConfigDB) GetConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
	)
}
