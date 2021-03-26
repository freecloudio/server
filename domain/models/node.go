package models

import (
	"time"
)

type NodeID string
type NodeType string
type NodeMimeType string

const (
	NodeTypeFile   NodeType = "FILE"
	NodeTypeFolder NodeType = "FOLDER"
)

type Node struct {
	ID      NodeID    `json:"id" fc_neo:",unique"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Name     string       `json:"name"`
	Size     int64        `json:"size"`
	MimeType NodeMimeType `json:"mime_type" fc_neo:",optional"`

	OwnerID   UserID    `json:"owner_id" fc_neo:"-"`
	Type      NodeType  `json:"type" fc_neo:"-"`
	ShareMode ShareMode `json:"share_mode" fc_neo:"-"`
	IsStarred bool      `json:"is_starred" fc_neo:"-"`
	Path      string    `json:"path" fc_neo:"-"`
	FullPath  string    `json:"full_path" fc_neo:"-"`

	PerspectiveUserID UserID `json:"-"`
}
