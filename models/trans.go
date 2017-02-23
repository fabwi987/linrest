package models

import (
	"database/sql"
	"strconv"
	"time"
)

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

func (db *DB) SumTransactionsByUser(userid int) int {

	stmt, err := db.Prepare("SELECT SUM(reward) AS TotalReward FROM test.trans WHERE iduser=?")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	rows, err := stmt.Query(userid)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var sum int
	for rows.Next() {
		err = rows.Scan(&sum)
		if err != nil {
			return 0
		}
	}

	if err = rows.Err(); err != nil {
		return 0
	}
	return sum
}

func (db *DB) CreateTransaction(rec int, user int, reward int) (sql.Result, error) {

	stmt, err := db.Prepare("INSERT trans SET idrec=?, iduser=?, reward=?, created=?, lastupdated=?")
	if err != nil {
		return nil, err
	}

	created := time.Now()

	res, err := stmt.Exec(rec, user, reward, created, created)
	if err != nil {
		return nil, err
	}

	inserteid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE trans SET url=? WHERE idtrans=" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	res, err = stmt.Exec("/trans/" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	return res, nil

}
