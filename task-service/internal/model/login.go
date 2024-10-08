package model

import "time"

var (
	SessionTime = 1 * time.Hour
)

type LoginCacheData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
