package models

import "time"

type Token string

type Session struct {
	Token      Token     `json:"token" fc_neo:",index"`
	UserID     UserID    `json:"user_id" fc_neo:"-"`
	ValidUntil time.Time `json:"valid_until"`
}
