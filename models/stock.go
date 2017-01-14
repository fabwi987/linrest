package models

//StockDataSaveFormat is the format for saving the stock data
type Stock struct {
	ID                 int     `json:"ID" bson:"ID"`
	Name               string  `json:"Name" bson:"Name"`
	Symbol             string  `json:"Symbol" bson:"Symbol"`
	LastTradePriceOnly float64 `json:"LastTradePriceOnly" bson:"LastTradePriceOnly"`
	Change             float64 `json:"Change" bson:"Change"`
	BuyPrice           float64 `json:"BuyPrice" bson:"BuyPrice"`
	NumberOfShares     float64 `json:"NumberOfShares" bson:"NumberOfShares"`
	Created            string  `json:"created" bson:"Created"`
	SalesPrice         float64 `json:"SalesPrice" bson:"SalesPrice"`
}

//GetStocks return all stocks from the database
func (db *DB) GetStocks() ([]*Stock, error) {

	rows, err := db.Query("SELECT * FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stcks := make([]*Stock, 0)
	for rows.Next() {
		tempStock := new(Stock)
		err = rows.Scan(&tempStock.ID, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares, &tempStock.SalesPrice, &tempStock.Name, &tempStock.LastTradePriceOnly)
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

	stmt, err := db.Prepare("SELECT * FROM stocks WHERE symbol = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()
	tempStock := new(Stock)

	for rows.Next() {
		err := rows.Scan(&tempStock.ID, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares, &tempStock.SalesPrice, &tempStock.Name, &tempStock.LastTradePriceOnly)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tempStock, nil
}
