package utils

import (
	"encoding/json"
	"net/http"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
)

//go:generate mockgen -source=response-writer.go -destination=response-writer_mock.go -package=utils
type IResponseWriter interface {
	SendSuccessResponse(http.ResponseWriter, int, string, any)
	SendErrorResponse(http.ResponseWriter, int, string, any)
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

func (re *ResponseWriter) SendErrorResponse(w http.ResponseWriter, status int, message string, reason any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := model.BaseResponse{
		StatusCode: status,
		Message:    message,
		Error:      reason,
	}

	json.NewEncoder(w).Encode(res)
}
