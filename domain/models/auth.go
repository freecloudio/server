package models

import "time"

type TokenValue string

type Token struct {
	Value      TokenValue `json:"value,omitempty" fc_neo:",index"`
	UserID     UserID     `json:"user_id,omitempty" fc_neo:"-"`
	ValidUntil time.Time  `json:"valid_until,omitempty"`
}
