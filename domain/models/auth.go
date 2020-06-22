package models

import "time"

type TokenValue string

type Token struct {
	Value      TokenValue `json:"value,omitempty"`
	UserID     UserID     `json:"user_id,omitempty"`
	ValidUntil time.Time  `json:"valid_until,omitempty"`
}
