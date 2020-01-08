package utils

import (
	models "DynamicAPI/model"
	"encoding/json"
	"net/http"
)

//SendError response error
func SendError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	//encode
	json.NewEncoder(w).Encode(error)
}

//SendSuccess response success
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//encode
	json.NewEncoder(w).Encode(data)
}
