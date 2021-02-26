package models

import "time"

type Token string

type Session struct {
	Token      Token     `json:"token,omitempty" fc_neo:",index"`
	UserID     UserID    `json:"user_id,omitempty" fc_neo:"-"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
