package main

import (
	"fmt"
	"net/http"
	"travelagency/api"
	"travelagency/repository"

	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
)

func main() {
	bannerFigure := figure.NewFigure("Travel Agency API", "", true)
	bannerFigure.Print()

	err := repository.EnsureDBExists()
	if err != nil {
		fmt.Println("error initializing database")
		return
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/health", getHealth())
	router.HandleFunc("/restart", getRestartDB())

	router.HandleFunc("/locations", api.RespondLocations)
	router.HandleFunc("/locations/{id}", api.RespondLocationDetails)

	router.HandleFunc("/holidays", api.RespondHolidays)
	router.HandleFunc("/holidays/{id}", api.RespondHolidayDetails)

	router.HandleFunc("/reservations", api.RespondReservations)
	router.HandleFunc("/reservations/{id}", api.RespondReservationDetails)

	server := &http.Server{Addr: "127.0.0.1:8080", Handler: router}
	server.ListenAndServe()
}

func getHealth() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		statusCode := http.StatusOK
		writer.WriteHeader(statusCode)
		writer.Write([]byte("{status:UP}"))
	}
}

func getRestartDB() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := repository.RestartDB()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("sucess"))
	}
}
