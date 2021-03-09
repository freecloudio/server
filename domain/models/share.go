package models

type ShareMode int

const (
	ShareModeNone ShareMode = iota // Only used in Node struct
	ShareModeRead
	ShareModeReadWrite
)

type Share struct {
	NodeID       NodeID    `json:"node_id"`
	OwnerID      UserID    `json:"owner_id"`
	SharedWithID UserID    `json:"shared_with_id"`
	Mode         ShareMode `json:"share_mode"`
}
