package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Datastore interface {
	GetStocks() ([]*Stock, error)
	GetSingleStock(symbol string) (*Stock, error)
	GetUsers() ([]*User, error)
	GetSingleUser(symbol int) (*User, error)
	GetRecommendations() ([]*Recommendation, error)
	GetRecommendationsByUser(symbol int) ([]*Recommendation, error)
}

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
