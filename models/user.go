package models

type User struct {
	ID    int    `json:"ID" bson:"ID"`
	Name  string `json:"Name" bson:"Name"`
	Phone string `json:"Phone" bson:"Phone"`
	Mail  string `json:"Mail" bson:"Mail"`
}

func (db *DB) GetUsers() ([]*User, error) {

	rows, err := db.Query("SELECT idusers, name, phone, mail FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usrs := make([]*User, 0)
	for rows.Next() {
		tempUser := new(User)
		err = rows.Scan(&tempUser.ID, &tempUser.Name, &tempUser.Phone, &tempUser.Mail)
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

	stmt, err := db.Prepare("SELECT idusers, name, phone, mail FROM users WHERE idusers = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()
	tempUser := new(User)

	for rows.Next() {
		err := rows.Scan(&tempUser.ID, &tempUser.Name, &tempUser.Phone, &tempUser.Mail)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tempUser, nil
}
