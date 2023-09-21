package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	uuid "github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

type DeliveryDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentDTO struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       uint   `json:"amount"`
	PaymentDT    uint   `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost uint   `json:"delivery_cost"`
	GoodsTotal   uint   `json:"goods_total"`
	CustomFee    uint   `json:"custom_fee"`
}

type ItemDTO struct {
	ChrtID      uint   `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       uint   `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        uint   `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  uint   `json:"total_price"`
	NmID        uint   `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      uint   `json:"status"`
}

type OrderDTO struct {
	OrderUID          string      `json:"order_uid"`
	TrackNumber       string      `json:"track_number"`
	Entry             string      `json:"entry"`
	Delivery          DeliveryDTO `json:"delivery"`
	Payment           PaymentDTO  `json:"payment"`
	Items             []ItemDTO   `json:"items"`
	Locale            string      `json:"locale"`
	InternalSignature string      `json:"internal_signature"`
	CustomerID        string      `json:"customer_id"`
	DeliveryService   string      `json:"delivery_service"`
	Shardkey          string      `json:"shardkey"`
	SmID              uint        `json:"sm_id"`
	DateCreated       string      `json:"date_created"`
	OofShard          string      `json:"oof_shard"`
}

const (
	clientID  = "sender"
	clusterID = "test-cluster"
	subject   = "order-channel"
	orderJSON = `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`
)

func main() {
	natsHost := os.Getenv("NATS_HOST")
	natsPort := os.Getenv("NATS_PORT")
	natsURI := fmt.Sprintf("nats://%s:%s", natsHost, natsPort)
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURI))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer nc.Close()
	log.Println("INFO nats: Successfully conected to nats-striaming")

	order := OrderDTO{}
	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		log.Fatal(err)
	}

	for {
		uuidObj, err := uuid.NewRandom()
		if err != nil {
			fmt.Printf("Error generating UUID: %v\n", err)
			return
		}
		uuidStr := uuidObj.String()
		order.OrderUID = uuidStr
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		err = nc.Publish(subject, jsonData)
		if err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
		fmt.Printf("INFO nats: Published new order with id -> %s\n", uuidStr)
		time.Sleep(5 * time.Second)
	}
}
