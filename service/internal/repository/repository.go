package repository

import (
	"database/sql"
	"service/internal/config"

	_ "github.com/lib/pq"
)

type Repository struct {
	cfg *config.ConfigDB
}

func New(cfg *config.ConfigDB) *Repository {
	return &Repository{cfg}
}

type OrderJSON []byte

func (r *Repository) GetOrders() ([]OrderJSON, error) {
	db, err := sql.Open("postgres", r.cfg.GetConnStr())
	if err != nil {
		return []OrderJSON{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT data FROM orders")
	if err != nil {
		return []OrderJSON{}, err
	}
	defer rows.Close()

	ordersJSON := []OrderJSON{}
	for rows.Next() {
		var o OrderJSON
		if err := rows.Scan(&o); err != nil {
			return []OrderJSON{}, err
		}
		ordersJSON = append(ordersJSON, o)
	}

	if err := rows.Err(); err != nil {
		return []OrderJSON{}, err
	}
	return ordersJSON, nil
}

func (r *Repository) InsertOrder(o OrderJSON) error {
	db, err := sql.Open("postgres", r.cfg.GetConnStr())
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO orders(data) VALUES($1)", o)
	if err != nil {
		return err
	}

	return nil
}
