package api

import "net/http"

type APIResponse struct {
	Status      int     `json:"status,omitempty"`
	Content     []byte  `json:"content,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
}

var (
	ContentTypeJSON = "application/json"

	ContentUnauthorized        = "Unauthorized\n"
	ContentOK                  = "OK\n"
	ContentInternalServerError = "Internal Server Error\n"
	ContentBadRequestError     = "Bad Request\n"
	ContentNotFoundError       = "Not Found\n"
)

func DefaultNotFoundError() APIResponse {
	return NotFoundError([]byte(ContentNotFoundError))
}

func NotFoundError(content []byte) APIResponse {
	return APIResponse{
		Status:  http.StatusNotFound,
		Content: content,
	}
}

func DefaultInternalServerError() APIResponse {
	return InternalServerError([]byte(ContentInternalServerError))
}

func InternalServerError(content []byte) APIResponse {
	return APIResponse{
		Status:  http.StatusInternalServerError,
		Content: content,
	}
}

func DefaultBadRequestError() APIResponse {
	return BadRequestError([]byte(ContentBadRequestError))
}

func BadRequestError(content []byte) APIResponse {
	return APIResponse{
		Status:  http.StatusBadRequest,
		Content: content,
	}
}

func DefaultOK() APIResponse {
	return OK([]byte(ContentOK))
}

func OKContentType(content []byte, contentType string) APIResponse {
	return APIResponse{
		Status:      http.StatusOK,
		Content:     content,
		ContentType: &contentType,
	}
}

func OK(content []byte) APIResponse {
	return APIResponse{
		Status:  http.StatusOK,
		Content: content,
	}
}
