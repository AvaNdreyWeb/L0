package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	order "service/internal/order"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	stan "github.com/nats-io/stan.go"
)

type DataJSON struct {
	Code int    `json:"code"`
	Id   string `json:"id"`
}

const (
	clientID    = "reciver"
	clusterID   = "test-cluster"
	subject     = "order-channel"
	queue       = "order-queue"
	durableName = "order-durable"
)

var cache = make(map[string]order.OrderDTO)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("INFO postgres: Successfully conected to database\n%s\n", connStr)
	rows, err := db.Query("SELECT data FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderJSON []byte
		if err := rows.Scan(&orderJSON); err != nil {
			log.Fatal(err)
		}

		newOrder := order.OrderDTO{}
		if err := json.Unmarshal(orderJSON, &newOrder); err != nil {
			log.Fatal(err)
		}

		cache[newOrder.OrderUID] = newOrder
		fmt.Printf("INFO cache: From postgres order %s\n", newOrder.OrderUID)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		var res []byte
		var err error
		if data, ok := cache[id]; ok {
			res, err = json.Marshal(data)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("INFO server: GET http://localhost:8080/ -> 200 OK")
		} else {
			data := struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Code:    404,
				Message: "Order not found",
			}
			res, err = json.Marshal(data)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("INFO server: GET http://localhost:8080/ -> 404 NOT FOUND")
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
	sub, err := nc.QueueSubscribe(subject, queue, func(msg *stan.Msg) {
		newOrder := order.OrderDTO{}
		if err := json.Unmarshal(msg.Data, &newOrder); err != nil {
			log.Fatal(err)
		}
		cache[newOrder.OrderUID] = newOrder
		dbOrder, err := json.Marshal(newOrder)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO orders(data) VALUES($1)", dbOrder)
		if err != nil {
			log.Fatalf("INSERT INTO orders(data) VALUES($1) error: %s", err)
		}

		fmt.Printf("INFO nats: Received order %s\n", newOrder.OrderUID)
	}, stan.DurableName(durableName))
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
