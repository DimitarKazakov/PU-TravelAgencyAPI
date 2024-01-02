package repository

import "database/sql"

type HolidaysEntity struct {
	ID         int64           `json:"id"`
	Title      string          `json:"title"`
	StartDate  string          `json:"startDate"`
	Duration   int             `json:"duration"`
	Price      float64         `json:"price"`
	FreeSlots  int             `json:"freeSlots"`
	LocationId int64           `json:"-"` //? used only to query the location entity from db
	Location   LocationsEntity `json:"location"`
}

func (hol *HolidaysRepo) Insert(entity HolidaysEntity) (*HolidaysEntity, error) {
	hol.mu.Lock()
	defer hol.mu.Unlock()

	resp, err := hol.db.Exec("INSERT INTO holidays VALUES(NULL,?,?,?,?,?,?);", entity.Title, entity.StartDate, entity.Duration, entity.Price, entity.FreeSlots, entity.LocationId)
	if err != nil {
		return nil, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}

	responseData := HolidaysEntity{
		ID:         id,
		Title:      entity.Title,
		StartDate:  entity.StartDate,
		Duration:   entity.Duration,
		Price:      entity.Price,
		FreeSlots:  entity.FreeSlots,
		LocationId: entity.LocationId,
	}

	locationEntity, err := hol.locationRepo.GetByID(responseData.LocationId)
	if err != nil {
		return nil, err
	}

	responseData.Location = *locationEntity
	return &entity, nil
}

func (hol *HolidaysRepo) Update(entity HolidaysEntity) (*HolidaysEntity, error) {
	hol.mu.Lock()
	defer hol.mu.Unlock()

	_, err := hol.db.Exec("UPDATE holidays SET title = ?, startDate = ?, duration = ?, price = ?, freeSlots = ?, locationId = ?  WHERE id = ?;", entity.Title, entity.StartDate, entity.Duration, entity.Price, entity.FreeSlots, entity.LocationId, entity.ID)
	if err != nil {
		return nil, err
	}

	locationEntity, err := hol.locationRepo.GetByID(entity.LocationId)
	if err != nil {
		return nil, err
	}

	entity.Location = *locationEntity
	return &entity, nil
}

func (hol *HolidaysRepo) GetAll(location string, startDate string, duration string) ([]HolidaysEntity, error) {
	query := "SELECT h.id, h.title, h.startDate, h.duration, h.price, h.freeSlots, h.locationId FROM holidays h JOIN locations l ON l.id = h.locationId WHERE 1=1"
	args := []interface{}{}

	if location != "" {
		query += " AND (l.city = ? OR l.country = ?) "
		args = append(args, location)
		args = append(args, location)
	}

	if startDate != "" {
		query += " AND h.startDate = ? "
		args = append(args, startDate)
	}

	if duration != "" {
		query += " AND h.duration = ? "
		args = append(args, duration)
	}

	query += ";"

	rows, err := hol.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []HolidaysEntity{}
	for rows.Next() {
		entity := HolidaysEntity{}
		err = rows.Scan(&entity.ID, &entity.Title, &entity.StartDate, &entity.Duration, &entity.Price, &entity.FreeSlots, &entity.LocationId)
		if err != nil {
			return nil, err
		}

		locationEntity, err := hol.locationRepo.GetByID(entity.LocationId)
		if err != nil {
			return nil, err
		}

		entity.Location = *locationEntity
		data = append(data, entity)
	}

	return data, nil
}

func (hol *HolidaysRepo) GetByID(id int64) (*HolidaysEntity, error) {
	row := hol.db.QueryRow("SELECT * FROM holidays WHERE id = ?;", id)

	var err error
	entity := HolidaysEntity{}
	if err = row.Scan(&entity.ID, &entity.Title, &entity.StartDate, &entity.Duration, &entity.Price, &entity.FreeSlots, &entity.LocationId); err == sql.ErrNoRows {
		return nil, err
	}

	locationEntity, err := hol.locationRepo.GetByID(entity.LocationId)
	if err != nil {
		return nil, err
	}

	entity.Location = *locationEntity
	return &entity, nil
}

func (hol *HolidaysRepo) Delete(id int64) error {
	hol.mu.Lock()
	defer hol.mu.Unlock()

	_, err := hol.db.Exec("DELETE FROM holidays WHERE id = ?;", id)
	return err
}
