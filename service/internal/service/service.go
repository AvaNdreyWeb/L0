package service

import (
	"encoding/json"
	"log"
	"service/internal/config"
	"service/internal/repository"
	"service/pkg/models"

	"github.com/nats-io/stan.go"
)

type Service struct {
	cfg   *config.ConfigNATS
	repo  *repository.Repository
	cache map[string]models.OrderModel
}

func New(repo *repository.Repository, cfg *config.ConfigNATS) *Service {
	cache := make(map[string]models.OrderModel)
	return &Service{cfg, repo, cache}
}

func (s *Service) Init() error {
	ordersJSON, err := s.repo.GetOrders()
	if err != nil {
		return err
	}

	for _, orderJSON := range ordersJSON {
		orderDTO := models.OrderModel{}
		if err := json.Unmarshal(orderJSON, &orderDTO); err != nil {
			return err
		}
		s.cache[orderDTO.OrderUID] = orderDTO
		log.Printf(
			"INFO: From DB {\"order_uid\": \"%s\", ...}",
			orderDTO.OrderUID,
		)
	}
	return nil
}

func (s *Service) GetOrderByID(id string) (models.OrderModel, bool) {
	orderDTO, ok := s.cache[id]
	return orderDTO, ok
}

func (s *Service) InsertOrder(msg *stan.Msg) {
	orderDTO := models.OrderModel{}
	if err := json.Unmarshal(msg.Data, &orderDTO); err != nil {
		log.Fatal(err)
	}
	s.cache[orderDTO.OrderUID] = orderDTO
	orderJSON, err := json.Marshal(orderDTO)
	if err != nil {
		log.Fatal(err)
	}
	s.repo.InsertOrder(orderJSON)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"INFO: From NATS {\"order_uid\": \"%s\", ...}",
		orderDTO.OrderUID,
	)
}

func (s *Service) Connect() (stan.Conn, error) {
	nc, err := stan.Connect(
		s.cfg.Cluster,
		s.cfg.Client,
		stan.NatsURL(s.cfg.GetConnStr()),
	)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

func (s *Service) Subscribe(nc stan.Conn) (stan.Subscription, error) {
	sub, err := nc.QueueSubscribe(
		s.cfg.Channel,
		s.cfg.Queue,
		s.InsertOrder,
		stan.DurableName(s.cfg.Durable),
	)
	if err != nil {
		return nil, err
	}
	return sub, nil
}
