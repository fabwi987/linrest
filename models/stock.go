package models

import (
	"database/sql"
	"strconv"
	"time"
)

//StockDataSaveFormat is the format for saving the stock data
type Stock struct {
	ID                 int       `json:"ID" bson:"ID"`
	Name               string    `json:"Name" bson:"Name"`
	Symbol             string    `json:"Symbol" bson:"Symbol"`
	LastTradePriceOnly float64   `json:"LastTradePriceOnly" bson:"LastTradePriceOnly"`
	Change             string    `json:"Change" bson:"Change"`
	BuyPrice           float64   `json:"BuyPrice" bson:"BuyPrice"`
	NumberOfShares     float64   `json:"NumberOfShares" bson:"NumberOfShares"`
	Created            time.Time `json:"created" bson:"Created"`
	Color              string    `json:"Color" bson:"Color"`
	SalesPrice         float64   `json:"SalesPrice" bson:"SalesPrice"`
	LastUpdated        time.Time `json:"LastUpdated" bson:"LastUpdated"`
	URL                string    `json:"URL" bson:"URL"`
}

//GetStocks return all stocks from the database
func (db *DB) GetStocks() ([]*Stock, error) {

	rows, err := db.Query("SELECT id, symbol, created, buyprice, numberofshare, salesprice, name, lasttradeprice, lastupdated, url FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stcks := make([]*Stock, 0)
	for rows.Next() {
		tempStock := new(Stock)
		err = rows.Scan(&tempStock.ID, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares, &tempStock.SalesPrice, &tempStock.Name, &tempStock.LastTradePriceOnly, &tempStock.LastUpdated, &tempStock.URL)
		if err != nil {
			return nil, err
		}

		stcks = append(stcks, tempStock)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stcks, nil
}

//GetSingleStock selects one specific stocks from the database
func (db *DB) GetSingleStock(symbol string) (*Stock, error) {

	stmt, err := db.Prepare("SELECT id, symbol, created, buyprice, numberofshare, salesprice, name, lasttradeprice, lastupdated, url FROM stocks WHERE symbol = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()
	tempStock := new(Stock)

	for rows.Next() {
		err := rows.Scan(&tempStock.ID, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares, &tempStock.SalesPrice, &tempStock.Name, &tempStock.LastTradePriceOnly, &tempStock.LastUpdated, &tempStock.URL)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tempStock, nil
}

func (db *DB) CreateStock(symbol string, createdorg time.Time, buyprice float64, numberofshares int, salesprice float64, name string, lasttradeprice float64) (sql.Result, error) {

	stmt, err := db.Prepare("INSERT stocks SET symbol=?, created=?, buyprice=?, numberofshare=?, salesprice=?, name=?, lasttradeprice=?, lastupdated=?")
	if err != nil {
		return nil, err
	}

	created := time.Now()

	res, err := stmt.Exec(symbol, createdorg, buyprice, numberofshares, salesprice, name, lasttradeprice, created)
	if err != nil {
		return nil, err
	}

	inserteid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE stocks SET url=? WHERE id=" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	res, err = stmt.Exec("/stocks/" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	return res, nil

}
