package models

type ShareMode string

const (
	ShareModeNone      ShareMode = "NONE" // Only used in Node struct
	ShareModeRead      ShareMode = "READ"
	ShareModeReadWrite ShareMode = "READ_WRITE"
)

type Share struct {
	NodeID       NodeID    `json:"node_id"`
	OwnerID      UserID    `json:"owner_id"`
	SharedWithID UserID    `json:"shared_with_id"`
	Mode         ShareMode `json:"share_mode"`
}
