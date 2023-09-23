package config

import (
	"log"
)

type Config struct {
	DB     *ConfigDB
	NATS   *ConfigNATS
	Server *ConfigServer
}

func New() *Config {
	// Initializing DB
	cfgDB := &ConfigDB{}
	if err := cfgDB.Getenv(); err != nil {
		log.Fatal(err)
	}

	// Initializing NATS
	cfgNATS := &ConfigNATS{}
	if err := cfgNATS.Getenv(); err != nil {
		log.Fatal(err)
	}

	// Initializing Server
	cfgServer := &ConfigServer{}
	if err := cfgServer.Getenv(); err != nil {
		log.Fatal(err)
	}

	return &Config{cfgDB, cfgNATS, cfgServer}
}
