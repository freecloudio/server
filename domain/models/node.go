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

	Size     int64        `json:"size"`
	MimeType NodeMimeType `json:"mime_type" fc_neo:",optional"`

	Name         string   `json:"name" fc_neo:"-"`
	OwnerID      UserID   `json:"owner_id" fc_neo:"-"`
	ParentNodeID *NodeID  `json:"parent_node_id" fc_neo:"-"`
	Type         NodeType `json:"type" fc_neo:"-"`
	IsStarred    bool     `json:"is_starred" fc_neo:"-"`
	Path         string   `json:"path" fc_neo:"-"`
	FullPath     string   `json:"full_path" fc_neo:"-"`

	PerspectiveUserID UserID `json:"-"`
}
