package fcerror

import (
	"encoding/json"
	"fmt"
)

type ErrorID int

type Error struct {
	ID          ErrorID `json:"id"`
	Description string  `json:"description"`
	Cause       error   `json:"cause"`
}

func (err Error) Error() string {
	return fmt.Sprintf("<%d> %s: %v", err.ID, err.Description, err.Cause)
}

func (err *Error) MarshalJSON() ([]byte, error) {
	type ShadowError Error
	type Default struct {
		*ShadowError
		CauseMsg string `json:"cause"`
	}
	type Nested struct {
		*ShadowError
		Cause Error `json:"cause"`
	}

	switch err.Cause.(type) {
	case Error:
		return json.Marshal(&Nested{
			ShadowError: (*ShadowError)(err),
			Cause:       err.Cause.(Error),
		})
	default:
		return json.Marshal(&Default{
			ShadowError: (*ShadowError)(err),
			CauseMsg:    err.Cause.Error(),
		})
	}
}

var errorDescriptions = map[ErrorID]string{}

func NewError(id ErrorID, cause error) *Error {
	return &Error{
		ID:          id,
		Description: errorDescriptions[id],
		Cause:       cause,
	}
}
