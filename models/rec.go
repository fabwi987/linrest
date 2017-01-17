package models

type Recommendation struct {
	ID   int    `json:"ID" bson:"ID"`
	Usr  *User  `json:"User" bson:"User"`
	Stck *Stock `json:"Stock" bson:"Stock"`
}

func (db *DB) GetRecommendations() ([]*Recommendation, error) {

	rows, err := db.Query("SELECT idrecs, idusr, idstock FROM recs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	poss := make([]*Recommendation, 0)
	var usr int
	var stc string
	for rows.Next() {
		ps := new(Recommendation)
		err = rows.Scan(&ps.ID, &usr, &stc)
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

	stmt, err := db.Prepare("SELECT idrecs, idusr, idstock FROM recs WHERE idusr = ?")
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
		err = rows.Scan(&ps.ID, &usr, &stc)
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
