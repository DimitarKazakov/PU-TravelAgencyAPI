package repository

import "database/sql"

type ReservationsEntity struct {
	ID          int64          `json:"id"`
	ContactName string         `json:"contactName"`
	PhoneNumber string         `json:"phoneNumber"`
	HolidayId   int64          `json:"-"` //? used only to query the holiday entity from db
	Holiday     HolidaysEntity `json:"holiday"`
}

func (res *ReservationsRepo) Insert(entity ReservationsEntity) (*ReservationsEntity, error) {
	res.mu.Lock()
	defer res.mu.Unlock()

	resp, err := res.db.Exec("INSERT INTO reservations VALUES(NULL,?,?,?);", entity.ContactName, entity.PhoneNumber, entity.HolidayId)
	if err != nil {
		return nil, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}

	responseData := ReservationsEntity{
		ID:          id,
		ContactName: entity.ContactName,
		PhoneNumber: entity.PhoneNumber,
		HolidayId:   entity.HolidayId,
	}

	holidayEntity, err := res.holidayRepo.GetByID(responseData.HolidayId)
	if err != nil {
		return nil, err
	}

	//? update free slots of holiday
	holidayEntity.FreeSlots -= 1
	_, err = res.holidayRepo.Update(*holidayEntity)
	if err != nil {
		return nil, err
	}

	responseData.Holiday = *holidayEntity
	return &responseData, nil
}

func (res *ReservationsRepo) Update(entity ReservationsEntity) (*ReservationsEntity, error) {
	res.mu.Lock()
	defer res.mu.Unlock()

	oldEntity, err := res.GetById(entity.ID)
	if err != nil {
		return nil, err
	}

	//? if we change the holiday we need to update free slots of old holiday
	if oldEntity.HolidayId != entity.HolidayId {
		oldHolidayEntity, err := res.holidayRepo.GetByID(oldEntity.HolidayId)
		if err != nil {
			return nil, err
		}

		oldHolidayEntity.FreeSlots += 1
		_, err = res.holidayRepo.Update(*oldHolidayEntity)
		if err != nil {
			return nil, err
		}
	}

	_, err = res.db.Exec("UPDATE reservations SET contactName = ?, phoneNumber = ?, holidayId = ? WHERE id = ?;", entity.ContactName, entity.PhoneNumber, entity.HolidayId, entity.ID)
	if err != nil {
		return nil, err
	}

	holidayEntity, err := res.holidayRepo.GetByID(entity.HolidayId)
	if err != nil {
		return nil, err
	}

	//? update free slots of holiday
	holidayEntity.FreeSlots -= 1
	_, err = res.holidayRepo.Update(*holidayEntity)
	if err != nil {
		return nil, err
	}

	entity.Holiday = *holidayEntity
	return &entity, nil
}

func (res *ReservationsRepo) GetAll() ([]ReservationsEntity, error) {
	rows, err := res.db.Query("SELECT * FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []ReservationsEntity{}
	for rows.Next() {
		entity := ReservationsEntity{}
		err = rows.Scan(&entity.ID, &entity.ContactName, &entity.PhoneNumber, &entity.HolidayId)
		if err != nil {
			return nil, err
		}

		holidayEntity, err := res.holidayRepo.GetByID(entity.HolidayId)
		if err != nil {
			return nil, err
		}

		entity.Holiday = *holidayEntity
		data = append(data, entity)
	}

	return data, nil
}

func (res *ReservationsRepo) GetById(id int64) (*ReservationsEntity, error) {
	row := res.db.QueryRow("SELECT * FROM reservations WHERE id = ?", id)

	var err error
	entity := ReservationsEntity{}
	if err = row.Scan(&entity.ID, &entity.ContactName, &entity.PhoneNumber, &entity.HolidayId); err == sql.ErrNoRows {
		return nil, err
	}

	holidayEntity, err := res.holidayRepo.GetByID(entity.HolidayId)
	if err != nil {
		return nil, err
	}

	entity.Holiday = *holidayEntity
	return &entity, nil
}

func (res *ReservationsRepo) Delete(id int64) error {
	res.mu.Lock()
	defer res.mu.Unlock()

	reservationEntity, err := res.GetById(id)
	if err != nil {
		return err
	}

	//? update free slots of holiday
	holidayEntity := reservationEntity.Holiday
	holidayEntity.FreeSlots += 1
	_, err = res.holidayRepo.Update(holidayEntity)
	if err != nil {
		return err
	}

	_, err = res.db.Exec("DELETE FROM reservations WHERE id = ?", id)
	return err
}
