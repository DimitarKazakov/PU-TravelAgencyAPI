package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"travelagency/repository"
)

type holidaysHandler struct {
	holidaysRepo *repository.HolidaysRepo
}

type holidayHandlerPostBody struct {
	Location  int64  `json:"location"`
	Title     string `json:"title"`
	StartDate string `json:"startDate"`
	Duration  int    `json:"duration"`
	Price     string `json:"price"`
	FreeSlots int    `json:"freeSlots"`
}

type holidayHandlerPutBody struct {
	ID        int64   `json:"id"`
	Location  int64   `json:"location"`
	Title     string  `json:"title"`
	StartDate string  `json:"startDate"`
	Duration  int     `json:"duration"`
	Price     float64 `json:"price"`
	FreeSlots int     `json:"freeSlots"`
}

func RespondHolidays(writer http.ResponseWriter, request *http.Request) {
	holidaysRepo, err := repository.NewHolidaysRepo(nil)
	if err != nil {
		response := InternalServerError([]byte("couldn't connect to database"))

		writer.WriteHeader(response.Status)
		writer.Write(response.Content)
		return
	}

	handler := holidaysHandler{
		holidaysRepo: holidaysRepo,
	}

	response := handler.respond(request)
	if response.ContentType != nil {
		writer.Header().Set("Content-Type", *response.ContentType)
	}

	writer.WriteHeader(response.Status)
	writer.Write(response.Content)
}

func (h *holidaysHandler) respond(request *http.Request) APIResponse {
	switch request.Method {

	case http.MethodGet:
		return h.handleGet(request)
	case http.MethodPost:
		return h.handlePost(request)
	case http.MethodPut:
		return h.handlePut(request)
	default:
		return InternalServerError([]byte("not implement\n"))
	}
}

func (h *holidaysHandler) handlePost(request *http.Request) APIResponse {
	var body holidayHandlerPostBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	parsedPrice, err := strconv.ParseFloat(body.Price, 64)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	entity, err := h.holidaysRepo.Insert(repository.HolidaysEntity{
		Title:      body.Title,
		StartDate:  body.StartDate,
		Duration:   body.Duration,
		Price:      parsedPrice,
		FreeSlots:  body.FreeSlots,
		LocationId: body.Location,
	})

	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *holidaysHandler) handlePut(request *http.Request) APIResponse {
	var body holidayHandlerPutBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	if body.ID == 0 {
		return BadRequestError([]byte("invalid holiday id"))
	}

	entity, err := h.holidaysRepo.Update(repository.HolidaysEntity{
		ID:         body.ID,
		Title:      body.Title,
		StartDate:  body.StartDate,
		Duration:   body.Duration,
		Price:      body.Price,
		FreeSlots:  body.FreeSlots,
		LocationId: body.Location,
	})

	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *holidaysHandler) handleGet(request *http.Request) APIResponse {
	locationFilter := request.URL.Query().Get("location")
	startDateFilter := request.URL.Query().Get("startDate")
	durationFilter := request.URL.Query().Get("duration")

	data, err := h.holidaysRepo.GetAll(locationFilter, startDateFilter, durationFilter)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(data)
	return OKContentType(jsonBody, ContentTypeJSON)
}
