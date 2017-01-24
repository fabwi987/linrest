package models

import "time"

type User struct {
	ID          int       `json:"ID" bson:"ID"`
	Name        string    `json:"Name" bson:"Name"`
	Phone       string    `json:"Phone" bson:"Phone"`
	Mail        string    `json:"Mail" bson:"Mail"`
	Created     time.Time `json:"Created" bson:"Created"`
	LastUpdated time.Time `json:"LastUpdated" bson:"LastUpdated"`
	URL         string    `json:"URL" bson:"URL"`
}

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
