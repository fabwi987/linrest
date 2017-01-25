package models

import (
	"database/sql"
	"strconv"
	"time"
)

type Meet struct {
	ID          int       `json:"ID" bson:"ID"`
	Location    string    `json:"Location" bson:"Location"`
	Date        time.Time `json:"Date" bson:"Date"`
	Text        string    `json:"Text" bson:"Text"`
	Created     time.Time `json:"Created" bson:"Created"`
	LastUpdated time.Time `json:"LastUpdated" bson:"LastUpdated"`
	URL         string    `json:"URL" bson:"URL"`
}

func (db *DB) GetMeets() ([]*Meet, error) {

	rows, err := db.Query("SELECT idmeet, location, date, text, created, lastupdated, url FROM meet")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poss := make([]*Meet, 0)
	for rows.Next() {
		ps := new(Meet)
		err = rows.Scan(&ps.ID, &ps.Location, &ps.Date, &ps.Text, &ps.Created, &ps.LastUpdated, &ps.URL)
		if err != nil {
			return nil, err
		}
		poss = append(poss, ps)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return poss, nil
}

func (db *DB) GetSingleMeet(id int) (*Meet, error) {

	stmt, err := db.Prepare("SELECT idmeet, location, date, text, created, lastupdated, url FROM meet WHERE idmeet=?")
	defer stmt.Close()
	rows, err := stmt.Query(id)
	defer rows.Close()
	tempUser := new(Meet)

	for rows.Next() {
		err := rows.Scan(&tempUser.ID, &tempUser.Location, &tempUser.Date, &tempUser.Text, &tempUser.Created, &tempUser.LastUpdated, &tempUser.URL)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tempUser, nil
}

func (db *DB) CreateMeet(location string, date time.Time, text string) (sql.Result, error) {

	stmt, err := db.Prepare("INSERT meet SET location=?, date=?, text=?, created=?, lastupdated=?")
	if err != nil {
		return nil, err
	}

	created := time.Now()

	res, err := stmt.Exec(location, date, text, created, created)
	if err != nil {
		return nil, err
	}

	inserteid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE meet SET url=? WHERE idmeet=" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	res, err = stmt.Exec("/meet/" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	return res, nil

}
