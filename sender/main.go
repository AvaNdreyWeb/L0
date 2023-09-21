package main

import (
	"fmt"
	"log"
	"os"
	"time"

	stan "github.com/nats-io/stan.go"
)

const (
	clientID       = "sender"
	clusterID      = "test-cluster"
	subject        = "order-channel"
	jsonDataToSend = `{"key1": "value1", "key2": "value2"}`
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

	jsonData := []byte(jsonDataToSend)

	for {
		err = nc.Publish(subject, jsonData)
		if err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
		fmt.Printf("INFO nats: Publish -> %s\n", jsonDataToSend)
		fmt.Printf("Published JSON data to channel '%s':\n%s\n", subject, jsonData)
		time.Sleep(5 * time.Second)
	}
}
