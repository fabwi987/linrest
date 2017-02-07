package models

import (
	"database/sql"
	"strconv"
	"time"
)

type Recommendation struct {
	ID          int       `json:"ID" bson:"ID"`
	Usr         *User     `json:"User" bson:"User"`
	Stck        *Stock    `json:"Stock" bson:"Stock"`
	Meet        *Meet     `json:"Meet" bson:"Meet"`
	Created     time.Time `json:"Created" bson:"Created"`
	LastUpdated time.Time `json:"LastUpdated" bson:"LastUpdated"`
	URL         string    `json:"URL" bson:"URL"`
}

type Recommendations []*Recommendation

func (slice Recommendations) Len() int { return len(slice) }
func (slice Recommendations) Less(i, j int) bool {
	return slice[i].Stck.Change > slice[j].Stck.Change
}
func (slice Recommendations) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func (db *DB) GetRecommendations() ([]*Recommendation, error) {

	rows, err := db.Query("SELECT idrecs, idusr, idstock, idmeet, created, lastupdated, url FROM recs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poss := make([]*Recommendation, 0)
	var usr int
	var stc string
	var mt int
	for rows.Next() {
		ps := new(Recommendation)
		err = rows.Scan(&ps.ID, &usr, &stc, &mt, &ps.Created, &ps.LastUpdated, &ps.URL)
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

		ps.Meet, err = db.GetSingleMeet(mt)
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

	stmt, err := db.Prepare("SELECT idrecs, idusr, idstock, idmeet, created, lastupdated, url FROM recs WHERE idusr = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poss := make([]*Recommendation, 0)
	var usr int
	var stc string
	var mt int
	for rows.Next() {
		ps := new(Recommendation)
		err = rows.Scan(&ps.ID, &usr, &stc, &mt, &ps.Created, &ps.LastUpdated, &ps.URL)
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

		ps.Meet, err = db.GetSingleMeet(mt)
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

func (db *DB) GetRecommendationsByMeet(symbol int) ([]*Recommendation, error) {

	stmt, err := db.Prepare("SELECT idrecs, idusr, idstock, idmeet, created, lastupdated, url FROM recs WHERE idmeet = ?")
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
	var mt int
	for rows.Next() {
		ps := new(Recommendation)
		err = rows.Scan(&ps.ID, &usr, &stc, &mt, &ps.Created, &ps.LastUpdated, &ps.URL)
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

		ps.Meet, err = db.GetSingleMeet(mt)
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

	inserteid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE recs SET url=? WHERE idrecs=" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	res, err = stmt.Exec("/recs/" + strconv.FormatInt(inserteid, 10))
	if err != nil {
		return nil, err
	}

	return res, nil

}
