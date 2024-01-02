package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"travelagency/repository"

	"github.com/gorilla/mux"
)

type locationDetailsHandler struct {
	locationsRepo *repository.LocationsRepo
}

func RespondLocationDetails(writer http.ResponseWriter, request *http.Request) {
	locationsRepo, err := repository.NewLocationsRepo(nil)
	if err != nil {
		response := InternalServerError([]byte("couldn't connect to database"))

		writer.WriteHeader(response.Status)
		writer.Write(response.Content)
		return
	}

	handler := locationDetailsHandler{
		locationsRepo: locationsRepo,
	}

	response := handler.respond(request)
	if response.ContentType != nil {
		writer.Header().Set("Content-Type", *response.ContentType)
	}

	writer.WriteHeader(response.Status)
	writer.Write(response.Content)
}

func (h *locationDetailsHandler) respond(request *http.Request) APIResponse {
	switch request.Method {

	case http.MethodGet:
		return h.handleGet(request)
	case http.MethodDelete:
		return h.handleDelete(request)
	default:
		return InternalServerError([]byte("not implemented\n"))
	}
}

func (h *locationDetailsHandler) handleGet(request *http.Request) APIResponse {
	vars := mux.Vars(request)
	idStr, exists := vars["id"]
	if !exists {
		return BadRequestError([]byte("id is empty"))
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	entity, err := h.locationsRepo.GetByID(id)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	if entity == nil {
		return DefaultNotFoundError()
	}

	jsonBody, _ := json.Marshal(entity)
	return OKContentType(jsonBody, ContentTypeJSON)
}

func (h *locationDetailsHandler) handleDelete(request *http.Request) APIResponse {
	vars := mux.Vars(request)
	idStr, exists := vars["id"]
	if !exists {
		return BadRequestError([]byte("id is empty"))
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	err = h.locationsRepo.Delete(id)
	if err != nil {
		return InternalServerError([]byte(err.Error()))
	}

	return OKContentType([]byte("true"), ContentTypeJSON)
}
