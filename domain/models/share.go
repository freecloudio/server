package models

type ShareMode string

const (
	ShareModeNone      ShareMode = ""
	ShareModeRead      ShareMode = "READ"
	ShareModeReadWrite ShareMode = "READ_WRITE"
)

type Share struct {
	NodeID       NodeID    `json:"node_id" fc_neo:"-"`
	SharedWithID UserID    `json:"shared_with_id" fc_neo:"-"`
	Mode         ShareMode `json:"share_mode"`
}
