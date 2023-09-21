package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	stan "github.com/nats-io/stan.go"
)

type DataJSON struct {
	Code int    `json:"code"`
	Id   string `json:"id"`
}

const (
	clientID  = "reciver"
	clusterID = "test-cluster"
	subject   = "order-channel"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("INFO postgres: Successfully conected to database")

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		log.Println("INFO server: GET http://localhost:8080/ -> 200 OK")

		data := DataJSON{200, id}
		res, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(res)
	}).Methods("GET")
	router.Use(mux.CORSMethodMiddleware(router))

	natsHost := os.Getenv("NATS_HOST")
	natsPort := os.Getenv("NATS_PORT")
	natsURI := fmt.Sprintf("nats://%s:%s", natsHost, natsPort)
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURI))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer nc.Close()
	log.Println("INFO nats: Successfully conected to nats-striaming")

	// Publish a message
	message := "Hello, NATS Streaming!"
	err = nc.Publish(subject, []byte(message))
	if err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}
	fmt.Printf("INFO nats: Publish -> %s\n", message)

	sub, err := nc.Subscribe(subject, func(msg *stan.Msg) {
		fmt.Printf("INFO nats: Received: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	defer sub.Unsubscribe()
	log.Println("INFO nats: Waiting for messages on order-channel")

	log.Println("INFO server: Starting HTTP server http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
