package api

import (
	"encoding/json"
	"net/http"
	"travelagency/repository"
)

type locationsHandler struct {
	locationsRepo *repository.LocationsRepo
}

type locationHandlerPostBody struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Number   string `json:"number"`
	Street   string `json:"street"`
	ImageUrl string `json:"imageUrl"`
}

type locationHandlerPutBody struct {
	ID       int64  `json:"id"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Number   string `json:"number"`
	Street   string `json:"street"`
	ImageUrl string `json:"imageUrl"`
}

func RespondLocations(writer http.ResponseWriter, request *http.Request) {
	locationsRepo, err := repository.NewLocationsRepo(nil)
	if err != nil {
		response := InternalServerError([]byte("couldn't connect to database"))

		writer.WriteHeader(response.Status)
		writer.Write(response.Content)
		return
	}

	handler := locationsHandler{
		locationsRepo: locationsRepo,
	}

	response := handler.respond(request)
	if response.ContentType != nil {
		writer.Header().Set("Content-Type", *response.ContentType)
	}

	writer.WriteHeader(response.Status)
	writer.Write(response.Content)
}

func (h *locationsHandler) respond(request *http.Request) APIResponse {
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

func (h *locationsHandler) handlePost(request *http.Request) APIResponse {
	var body locationHandlerPostBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	entity, err := h.locationsRepo.Insert(repository.LocationsEntity{
		City:     body.City,
		Country:  body.Country,
		Number:   body.Number,
		Street:   body.Street,
		ImageUrl: body.ImageUrl,
	})

	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *locationsHandler) handlePut(request *http.Request) APIResponse {
	var body locationHandlerPutBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	if body.ID == 0 {
		return BadRequestError([]byte("invalid location id"))
	}

	entity, err := h.locationsRepo.Update(repository.LocationsEntity{
		ID:       body.ID,
		City:     body.City,
		Country:  body.Country,
		Number:   body.Number,
		Street:   body.Street,
		ImageUrl: body.ImageUrl,
	})

	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *locationsHandler) handleGet(request *http.Request) APIResponse {
	data, err := h.locationsRepo.GetAll()
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	jsonBody, _ := json.Marshal(data)
	return OKContentType(jsonBody, ContentTypeJSON)
}
