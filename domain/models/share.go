package models

type ShareMode string

const (
	ShareModeNone      ShareMode = "NONE" // Only used in Node struct
	ShareModeRead      ShareMode = "READ"
	ShareModeReadWrite ShareMode = "READ_WRITE"
)

type Share struct {
	NodeID       NodeID    `json:"node_id" neo_fc:"-"`
	SharedWithID UserID    `json:"shared_with_id" fc_neo:"-"`
	Mode         ShareMode `json:"share_mode"`
	NameOverride string    `json:"name_override" neo_fc:"-"`
}
