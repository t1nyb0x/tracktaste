package handler

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func success(w http.ResponseWriter, result interface{}) {
	writeJSON(w, http.StatusOK, successResponse{Status: http.StatusOK, Result: result})
}

func badRequest(w http.ResponseWriter, message, code string) {
	writeJSON(w, http.StatusBadRequest, errorResponse{Status: http.StatusBadRequest, Message: message, Code: code})
}

func notFound(w http.ResponseWriter, message, code string) {
	writeJSON(w, http.StatusNotFound, errorResponse{Status: http.StatusNotFound, Message: message, Code: code})
}

func serviceUnavailable(w http.ResponseWriter, message, code string) {
	writeJSON(w, http.StatusServiceUnavailable, errorResponse{Status: http.StatusServiceUnavailable, Message: message, Code: code})
}
