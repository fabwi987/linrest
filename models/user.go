package models

import (
	"database/sql"
	"strconv"
	"time"
)

type User struct {
	ID          int       `json:"ID" bson:"ID"`
	Name        string    `json:"Name" bson:"Name"`
	Phone       string    `json:"Phone" bson:"Phone"`
	Mail        string    `json:"Mail" bson:"Mail"`
	Score       int       `json:"Score" bson:"Score"`
	Created     time.Time `json:"Created" bson:"Created"`
	LastUpdated time.Time `json:"LastUpdated" bson:"LastUpdated"`
	URL         string    `json:"URL" bson:"URL"`
}

type Users []*User

func (slice Users) Len() int { return len(slice) }
func (slice Users) Less(i, j int) bool {
	return slice[i].Score > slice[j].Score
}
func (slice Users) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func (db *DB) GetUsers() ([]*User, error) {

	rows, err := db.Query("SELECT idusers, name, phone, mail, created, lastupdated, url FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usrs := make([]*User, 0)
	for rows.Next() {
		tempUser := new(User)
		err = rows.Scan(&tempUser.ID, &tempUser.Name, &tempUser.Phone, &tempUser.Mail, &tempUser.Created, &tempUser.LastUpdated, &tempUser.URL)
		if err != nil {
			return nil, err
		}

		usrs = append(usrs, tempUser)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usrs, nil
}

func (db *DB) GetSingleUser(symbol int) (*User, error) {

	stmt, err := db.Prepare("SELECT idusers, name, phone, mail, created, lastupdated, url FROM users WHERE idusers = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()
	tempUser := new(User)

	for rows.Next() {
		err := rows.Scan(&tempUser.ID, &tempUser.Name, &tempUser.Phone, &tempUser.Mail, &tempUser.Created, &tempUser.LastUpdated, &tempUser.URL)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tempUser, nil
}

func (db *DB) GetUsersLeaderboard() ([]*User, error) {

	rows, err := db.Query("SELECT idusers, name, phone, mail, created, lastupdated, url FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//usrs := make(Users, 0)
	var usrs Users
	var iduser int
	for rows.Next() {
		tempUser := new(User)
		err = rows.Scan(&tempUser.ID, &tempUser.Name, &tempUser.Phone, &tempUser.Mail, &tempUser.Created, &tempUser.LastUpdated, &tempUser.URL)
		if err != nil {
			return nil, err
		}

		iduser = tempUser.ID
		tempUser.Score, err = db.SumTransactionsByUser(iduser)
		if err != nil {
			return nil, err
		}

		usrs = append(usrs, tempUser)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usrs, nil
}

func (db *DB) CreateUser(name string, phone string, mail string) (sql.Result, error) {

	stmt, err := db.Prepare("INSERT users SET name=?, phone=?, mail=?, created=?, lastupdated=?")
	if err != nil {
		return nil, err
	}

	created := time.Now()

	res, err := stmt.Exec(name, phone, mail, created, created)
	if err != nil {
		return nil, err
	}

	inserteid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE users SET url=? WHERE idusers=" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	res, err = stmt.Exec("/users/" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	return res, nil

}
