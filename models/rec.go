package models

import "time"
import "database/sql"

type Recommendation struct {
	ID          int       `json:"ID" bson:"ID"`
	Usr         *User     `json:"User" bson:"User"`
	Stck        *Stock    `json:"Stock" bson:"Stock"`
	Created     time.Time `json:"Created" bson:"Created"`
	LastUpdated time.Time `json:"LastUpdated" bson:"LastUpdated"`
}

func (db *DB) GetRecommendations() ([]*Recommendation, error) {

	rows, err := db.Query("SELECT idrecs, idusr, idstock, created, lastupdated FROM recs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poss := make([]*Recommendation, 0)
	var usr int
	var stc string
	for rows.Next() {
		ps := new(Recommendation)
		err = rows.Scan(&ps.ID, &usr, &stc, &ps.Created, &ps.LastUpdated)
		if err != nil {
			return nil, err
		}

		ps.Usr, err = db.GetSingleUser(usr)
		if err != nil {
			return nil, err
		}

		ps.Stck, err = db.GetSingleStock(stc)
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

func (db *DB) GetRecommendationsByUser(symbol int) ([]*Recommendation, error) {

	stmt, err := db.Prepare("SELECT idrecs, idusr, idstock, created, lastupdated FROM recs WHERE idusr = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poss := make([]*Recommendation, 0)
	var usr int
	var stc string
	for rows.Next() {
		ps := new(Recommendation)
		err = rows.Scan(&ps.ID, &usr, &stc, &ps.Created, &ps.LastUpdated)
		if err != nil {
			return nil, err
		}

		ps.Usr, err = db.GetSingleUser(usr)
		if err != nil {
			return nil, err
		}

		ps.Stck, err = db.GetSingleStock(stc)
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

func (db *DB) CreateRecommendation(symbol string, user int, meet int) (sql.Result, error) {

	stmt, err := db.Prepare("INSERT recs SET idusr=?, idstock=?, idmeet=?, created=?, lastupdated=?")
	if err != nil {
		return nil, err
	}

	created := time.Now()

	res, err := stmt.Exec(user, symbol, meet, created, created)
	if err != nil {
		return nil, err
	}

	return res, nil

}
