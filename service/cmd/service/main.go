package main

import (
	"log"
	"service/internal/config"
	"service/internal/repository"
	"service/internal/server"
	"service/internal/service"
)

func main() {
	cfg := config.New()

	repo := repository.New(cfg.DB)

	srvc := service.New(repo, cfg.NATS)
	err := srvc.Init()
	if err != nil {
		log.Fatal(err)
	}
	nc, err := srvc.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	sub, err := srvc.Subscribe(nc)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	server.New(srvc, cfg.Server).Run()
}
