package httpapi

import (
	"encoding/json"
	"net/http"
)

const (
	calculatePath            = "/api/calculate"
	jsonMediaType            = "application/json"
	codeNotFound             = "NOT_FOUND"
	codeMethodNotAllowed     = "METHOD_NOT_ALLOWED"
	codeUnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	codeInvalidRequest       = "INVALID_REQUEST"
	codeRequestTooLarge      = "REQUEST_TOO_LARGE"
)

type calculationResponse struct {
	Result string `json:"result"`
}
type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type errorResponse struct {
	Error errorBody `json:"error"`
}

type requestError struct {
	status  int
	code    string
	message string
}

func (err *requestError) write(w http.ResponseWriter) {
	if err.status == http.StatusMethodNotAllowed {
		w.Header().Set("Allow", http.MethodPost)
	}
	fail(w, err.status, err.code, err.message)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", jsonMediaType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func fail(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorBody{Code: code, Message: message}})
}
