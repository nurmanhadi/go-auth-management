package dto

import "time"

type EventUserPayload struct {
	Event     string        `json:"event"`
	Timestamp time.Time     `json:"timestamp"`
	Data      EventUserData `json:"data"`
}
type EventUserData struct {
	UserId        string    `json:"user_id"`
	Username      string    `json:"username"`
	Registered_at time.Time `json:"registered_at"`
}
