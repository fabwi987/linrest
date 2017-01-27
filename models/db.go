package models

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Datastore interface {
	GetStocks() ([]*Stock, error)
	GetSingleStock(symbol string) (*Stock, error)
	CreateStock(symbol string, createdorg time.Time, buyprice float64, numberofshares int, salesprice float64, name string, lasttradeprice float64) (sql.Result, error)

	GetUsers() ([]*User, error)
	GetSingleUser(symbol int) (*User, error)
	CreateUser(name string, phone string, mail string) (sql.Result, error)

	GetRecommendations() ([]*Recommendation, error)
	GetRecommendationsByUser(symbol int) ([]*Recommendation, error)
	CreateRecommendation(symbol string, user int, meet int) (sql.Result, error)

	GetMeets() ([]*Meet, error)
	GetSingleMeet(id int) (*Meet, error)
	CreateMeet(location string, date time.Time, text string) (sql.Result, error)
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
