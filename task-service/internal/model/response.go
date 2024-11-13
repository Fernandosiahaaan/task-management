package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateResponseHttp(w http.ResponseWriter, r *http.Request, statusCode int, response Response) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	if response.Error {
		fmt.Printf("❌  [%s] uri = '%s'; status code = %d; message = %s\n", r.Method, r.RequestURI, statusCode, response.Message)
		return
	}
	fmt.Printf("✅  [%s] uri = '%s'; status code = %d; message = %s\n", r.Method, r.RequestURI, statusCode, response.Message)
}

func ConvertResponseToStr(response Response) (string, error) {
	// Konversi struct menjadi JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	// Mengubah byte array menjadi string JSON
	return string(jsonData), nil
}
