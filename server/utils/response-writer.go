package utils

import (
	"encoding/json"
	"net/http"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
)

type IResponseWriter interface {
	SendSuccessResponse(w http.ResponseWriter, status int, message string, data any)
	SendErrorResponse(w http.ResponseWriter, status int, message string)
}

type ResponseWriter struct{}

func CreateNewResponseWriter() IResponseWriter {
	return &ResponseWriter{}
}

func (rw *ResponseWriter) SendSuccessResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := model.BaseResponse{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}

	json.NewEncoder(w).Encode(res)
}

func (re *ResponseWriter) SendErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := model.BaseResponse{
		StatusCode: status,
		Message:    message,
	}

	json.NewEncoder(w).Encode(res)
}
