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
	Id    string `json:"user_id"`
	Token string `json:"token"`
}

func CreateResponseHttp(w http.ResponseWriter, r *http.Request, statusCode int, response ResponseHttp) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	if response.Error {
		fmt.Printf("❌  [%s] uri = '%s'; status code = %d; message = %s\n", r.Method, r.RequestURI, statusCode, response.Message)
		return
	}
	fmt.Printf("✅  [%s] uri = '%s'; status code = %d; message = %s\n", r.Method, r.RequestURI, statusCode, response.Message)
}
