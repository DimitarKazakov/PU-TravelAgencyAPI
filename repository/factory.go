package repository

import (
	"database/sql"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: add prepared statements / defer stmt.Close() when you call them

var (
	databaseFile     string = "sqlite.db"
	connectionString string = "file:" + databaseFile + "?_foreign_keys=1"
)

type LocationsRepo struct {
	mu sync.Mutex
	db *sql.DB
}

type HolidaysRepo struct {
	mu sync.Mutex
	db *sql.DB

	locationRepo *LocationsRepo
}

type ReservationsRepo struct {
	mu sync.Mutex
	db *sql.DB

	holidayRepo *HolidaysRepo
}

func EnsureDBExists() error {
	if _, err := os.Stat(databaseFile); err == nil {
		return nil
	}

	return CreateDB()
}

func RestartDB() error {
	if _, err := os.Stat(databaseFile); err == nil {
		err := os.Remove(databaseFile)
		if err != nil {
			return err
		}
	}

	return CreateDB()
}

func CreateDB() error {
	_, err := os.Create(databaseFile)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return err
	}

	_, err = NewLocationsRepo(db)
	if err != nil {
		return err
	}

	_, err = NewHolidaysRepo(db)
	if err != nil {
		return err
	}

	_, err = NewReservationsRepo(db)
	if err != nil {
		return err
	}

	return nil
}

func NewLocationsRepo(db *sql.DB) (*LocationsRepo, error) {
	var err error
	if db == nil {
		db, err = sql.Open("sqlite3", connectionString)
		if err != nil {
			return nil, err
		}
	}

	createStatement := `
	CREATE TABLE IF NOT EXISTS locations (
		id INTEGER NOT NULL PRIMARY KEY,
		street TEXT NOT NULL,
		number TEXT NOT NULL,
		city TEXT NOT NULL,
		country TEXT NOT NULL,
		imageUrl TEXT NOT NULL
	);`

	if _, err := db.Exec(createStatement); err != nil {
		return nil, err
	}

	return &LocationsRepo{
		db: db,
	}, nil
}

func NewHolidaysRepo(db *sql.DB) (*HolidaysRepo, error) {
	var err error
	if db == nil {
		db, err = sql.Open("sqlite3", connectionString)
		if err != nil {
			return nil, err
		}
	}

	createStatement := `
	CREATE TABLE IF NOT EXISTS holidays (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		startDate TEXT NOT NULL,
		duration INTEGER NOT NULL,
		price REAL NOT NULL,
		freeSlots INTEGER NOT NULL,
		locationId INTEGER NOT NULL,
		FOREIGN KEY(locationId) REFERENCES locations(id)
	);`

	if _, err := db.Exec(createStatement); err != nil {
		return nil, err
	}

	locationRepo, err := NewLocationsRepo(db)
	if err != nil {
		return nil, err
	}

	return &HolidaysRepo{
		db:           db,
		locationRepo: locationRepo,
	}, nil
}

func NewReservationsRepo(db *sql.DB) (*ReservationsRepo, error) {
	var err error
	if db == nil {
		db, err = sql.Open("sqlite3", connectionString)
		if err != nil {
			return nil, err
		}
	}

	createStatement := `
	CREATE TABLE IF NOT EXISTS reservations (
		id INTEGER NOT NULL PRIMARY KEY,
		contactName TEXT NOT NULL,
		phoneNumber TEXT NOT NULL,
		holidayId INTEGER NOT NULL,
		FOREIGN KEY(holidayId) REFERENCES holidays(id)
	);`

	if _, err := db.Exec(createStatement); err != nil {
		return nil, err
	}

	holidayRepo, err := NewHolidaysRepo(db)
	if err != nil {
		return nil, err
	}

	return &ReservationsRepo{
		db:          db,
		holidayRepo: holidayRepo,
	}, nil
}
