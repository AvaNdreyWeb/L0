package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
		log.Println("INFO server: GET http://localhost:8080/ -> 200 OK")
		w.WriteHeader(http.StatusOK)
	})

	log.Println("INFO server: Starting HTTP server http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
