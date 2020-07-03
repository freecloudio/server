package fcerror

import "fmt"

type ErrorID int

type Error struct {
	ID          ErrorID `json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	Cause       error   `json:"cause,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("<%d> %s: %v", err.ID, err.Description, err.Cause)
}

var errorDescriptions = map[ErrorID]string{}

func NewError(id ErrorID, cause error) *Error {
	return &Error{
		ID:          id,
		Description: errorDescriptions[id],
		Cause:       cause,
	}
}
