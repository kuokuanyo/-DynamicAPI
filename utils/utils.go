package utils

import (
	models "DynamicAPI/model"
	"encoding/json"
	"net/http"
	"strings"
)

//SendError response error
func SendError(w http.ResponseWriter, status int, message models.Error, err error) {
	w.WriteHeader(status)
	//encode
	json.NewEncoder(w).Encode(message)
	json.NewEncoder(w).Encode(err)
}

//SendSuccess response success
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//encode
	json.NewEncoder(w).Encode(data)
}

//duplicate 檢查是否有重複字串
func Duplicate(col []string, s string) bool {
	b := true
	for x := range col {
		if strings.Contains(col[x], s) {
			b = false
			break
		} else {
			continue
		}
	}
	return b
}
