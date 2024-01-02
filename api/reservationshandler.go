package api

import (
	"encoding/json"
	"net/http"
	"travelagency/repository"
)

type reservationsHandler struct {
	reservationsRepo *repository.ReservationsRepo
}

type reservationsHandlerPostBody struct {
	ContactName string `json:"contactName"`
	PhoneNumber string `json:"phoneNumber"`
	Holiday     int64  `json:"holiday"`
}

type reservationsHandlerPutBody struct {
	ID          int64  `json:"id"`
	ContactName string `json:"contactName"`
	PhoneNumber string `json:"phoneNumber"`
	Holiday     int64  `json:"holiday"`
}

func RespondReservations(writer http.ResponseWriter, request *http.Request) {
	reservationsRepo, err := repository.NewReservationsRepo(nil)
	if err != nil {
		response := InternalServerError([]byte("couldn't connect to database"))

		writer.WriteHeader(response.Status)
		writer.Write(response.Content)
		return
	}

	handler := reservationsHandler{
		reservationsRepo: reservationsRepo,
	}

	response := handler.respond(request)
	if response.ContentType != nil {
		writer.Header().Set("Content-Type", *response.ContentType)
	}

	writer.WriteHeader(response.Status)
	writer.Write(response.Content)
}

func (h *reservationsHandler) respond(request *http.Request) APIResponse {
	switch request.Method {

	case http.MethodGet:
		return h.handleGet(request)
	case http.MethodPost:
		return h.handlePost(request)
	case http.MethodPut:
		return h.handlePut(request)
	default:
		return InternalServerError([]byte("not implemented\n"))
	}
}

func (h *reservationsHandler) handlePost(request *http.Request) APIResponse {
	var body reservationsHandlerPostBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	entity, err := h.reservationsRepo.Insert(repository.ReservationsEntity{
		ContactName: body.ContactName,
		PhoneNumber: body.PhoneNumber,
		HolidayId:   body.Holiday,
	})

	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *reservationsHandler) handlePut(request *http.Request) APIResponse {
	var body reservationsHandlerPutBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	entity, err := h.reservationsRepo.Update(repository.ReservationsEntity{
		ID:          body.ID,
		ContactName: body.ContactName,
		PhoneNumber: body.PhoneNumber,
		HolidayId:   body.Holiday,
	})

	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *reservationsHandler) handleGet(request *http.Request) APIResponse {
	data, err := h.reservationsRepo.GetAll()
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(data)
	return OKContentType(jsonBody, ContentTypeJSON)
}
