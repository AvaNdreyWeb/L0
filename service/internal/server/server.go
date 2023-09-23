package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"service/internal/config"
	"service/internal/service"

	"github.com/gorilla/mux"
)

type Server struct {
	cfg  *config.ConfigServer
	srvc *service.Service
}

func New(srvc *service.Service, cfg *config.ConfigServer) *Server {
	return &Server{cfg, srvc}
}

func (s *Server) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", s.getByIDHandler).Methods("GET")
	router.Use(mux.CORSMethodMiddleware(router))

	addr := s.cfg.GetAddr()
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) getByIDHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	orderDTO, ok := s.srvc.GetOrderByID(id)

	var res []byte
	var err error
	if ok {
		res, err = json.Marshal(orderDTO)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(
			"INFO: GET %s/?id=%s -> 200 OK",
			s.cfg.GetAddr(),
			id,
		)
	} else {
		notFoundDTO := struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    404,
			Message: fmt.Sprintf("Заказ с uid: '%s' не найден", id),
		}
		res, err = json.Marshal(notFoundDTO)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(
			"INFO: GET %s/?id=%s -> 404 NOT FOUND",
			s.cfg.GetAddr(),
			id,
		)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(res)
}
