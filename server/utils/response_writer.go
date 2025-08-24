package utils

import (
	"encoding/json"
	"net/http"

	"github.com/shashikedissanayake/web-page-analyzer/server/model"
)

//go:generate mockgen -source=response_writer.go -destination=response_writer_mock.go -package=utils
type IResponseWriter interface {
	SendSuccessResponse(http.ResponseWriter, int, string, any)
	SendErrorResponse(http.ResponseWriter, int, string, any)
}

var (
	SERVER_FAILURE_RESPONSE = &model.BaseResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to parse response",
	}
)

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

	resJson, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SERVER_FAILURE_RESPONSE)
	}

	w.WriteHeader(status)
	w.Write(resJson)
}

func (re *ResponseWriter) SendErrorResponse(w http.ResponseWriter, status int, message string, reason any) {
	w.Header().Set("Content-Type", "application/json")

	res := model.BaseResponse{
		StatusCode: status,
		Message:    message,
		Error:      reason,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SERVER_FAILURE_RESPONSE)
	}

	w.WriteHeader(status)
	w.Write(resJson)
}
