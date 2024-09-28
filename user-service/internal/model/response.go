package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseHttp struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginData struct {
	Token string `json:"token"`
}

func CreateResponseHttp(w http.ResponseWriter, statusCode int, response ResponseHttp) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	if response.Error {
		fmt.Printf("❌ status code = %d; message = %s\n", statusCode, response.Message)
		return
	}
	fmt.Printf("✔️  status code = %d; message = %s\n", statusCode, response.Message)
}
