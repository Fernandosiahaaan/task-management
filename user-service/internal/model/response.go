package model

type ResponseHttp struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginData struct {
	Token string `json:"token"`
}