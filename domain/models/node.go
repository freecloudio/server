package models

import (
	"time"
)

type NodeType int
type NodeID int64
type NodeMimeType string

const (
	NodeTypeFile NodeType = iota
	NodeTypeFolder
)

type Node struct {
	ID      NodeID    `json:"id" fc_neo:"-"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Name     string       `json:"name"`
	Type     NodeType     `json:"type"`
	OwnerID  UserID       `json:"owner_id"`
	Size     int64        `json:"size"`
	MimeType NodeMimeType `json:"mime_type" fc_neo:",optional"`

	PerspectiveUserID UserID    `json:"-"`
	ShareMode         ShareMode `json:"share_mode" fc_neo:"-"`
	IsStarred         bool      `json:"is_starred" fc_neo:"-"`
	Path              string    `json:"path" fc_neo:"-"`
	FullPath          string    `json:"full_path" fc_neo:"-"`
}
