package handlers

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}




func RespondSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIResponse{
        Status:  status,
        Message: message,
        Data:    data,
    })
}

func RespondError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIResponse{
        Status:  status,
        Message: message,
        Error:   message,
    })
}
