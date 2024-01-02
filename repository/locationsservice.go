package repository

import "database/sql"

type LocationsEntity struct {
	ID       int64  `json:"id"`
	Street   string `json:"street"`
	Number   string `json:"number"`
	City     string `json:"city"`
	Country  string `json:"country"`
	ImageUrl string `json:"imageUrl"`
}

func (loc *LocationsRepo) Insert(entity LocationsEntity) (*LocationsEntity, error) {
	loc.mu.Lock()
	defer loc.mu.Unlock()

	resp, err := loc.db.Exec("INSERT INTO locations VALUES(NULL,?,?,?,?,?);", entity.Street, entity.Number, entity.City, entity.Country, entity.ImageUrl)
	if err != nil {
		return nil, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &LocationsEntity{
		ID:       id,
		Street:   entity.Street,
		Number:   entity.Number,
		City:     entity.City,
		Country:  entity.Country,
		ImageUrl: entity.ImageUrl,
	}, nil
}

func (loc *LocationsRepo) Update(entity LocationsEntity) (*LocationsEntity, error) {
	loc.mu.Lock()
	defer loc.mu.Unlock()

	_, err := loc.db.Exec("UPDATE locations SET street=?, number=?, city=?, country=?, imageUrl=? WHERE id = ?;", entity.Street, entity.Number, entity.City, entity.Country, entity.ImageUrl, entity.ID)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (loc *LocationsRepo) GetAll() ([]LocationsEntity, error) {
	rows, err := loc.db.Query("SELECT * FROM locations;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []LocationsEntity{}
	for rows.Next() {
		entity := LocationsEntity{}
		err = rows.Scan(&entity.ID, &entity.Street, &entity.Number, &entity.City, &entity.Country, &entity.ImageUrl)
		if err != nil {
			return nil, err
		}

		data = append(data, entity)
	}

	return data, nil
}

func (loc *LocationsRepo) GetByID(id int64) (*LocationsEntity, error) {
	row := loc.db.QueryRow("SELECT * FROM locations WHERE id = ?;", id)

	var err error
	entity := LocationsEntity{}

	if err = row.Scan(&entity.ID, &entity.Street, &entity.Number, &entity.City, &entity.Country, &entity.ImageUrl); err == sql.ErrNoRows {
		return nil, err
	}

	return &entity, nil
}

func (loc *LocationsRepo) Delete(id int64) error {
	loc.mu.Lock()
	defer loc.mu.Unlock()

	_, err := loc.db.Exec("DELETE FROM locations WHERE id = ?;", id)
	return err
}
