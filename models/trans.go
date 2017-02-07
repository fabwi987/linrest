package models

import "time"

type Transaction struct {
	ID             int       `json:"ID" bson:"ID"`
	Recommendation int       `json:"Recommendation" bson:"Recommendation"`
	User           int       `json:"User" bson:"User"`
	Reward         string    `json:"Reward" bson:"Reward"`
	Created        time.Time `json:"Created" bson:"Created"`
	LastUpdated    time.Time `json:"LastUpdated" bson:"LastUpdated"`
	URL            string    `json:"URL" bson:"URL"`
}

func (db *DB) GetTransactionsByUser(userid int) ([]*Transaction, error) {

	stmt, err := db.Prepare("SELECT idtrans, idrec, iduser, reward, created, lastupdated, url FROM trans WHERE iduser=?")
	defer stmt.Close()
	rows, err := stmt.Query(userid)
	defer rows.Close()

	trans := make([]*Transaction, 0)
	for rows.Next() {
		ps := new(Transaction)
		err = rows.Scan(&ps.ID, &ps.Recommendation, &ps.User, &ps.Reward, &ps.Created, &ps.LastUpdated, &ps.URL)
		if err != nil {
			return nil, err
		}
		trans = append(trans, ps)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trans, nil
}

func (db *DB) SumTransactionsByUser(userid int) (*int, error) {

	stmt, err := db.Prepare("SELECT SUM(reward) AS TotalReward FROM test.trans WHERE iduser=?")
	defer stmt.Close()
	rows, err := stmt.Query(userid)
	defer rows.Close()

	var sum *int
	for rows.Next() {
		err = rows.Scan(&sum)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sum, nil
}
