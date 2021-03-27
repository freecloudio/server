package fcerror

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type ErrorID int

type Error struct {
	ID          ErrorID `json:"id"`
	Description string  `json:"description"`
	Cause       error   `json:"cause"`
	File        string  `json:"file,omitempty"`
	Line        int     `json:"line,omitempty"`
	Function    string  `json:"function,omitempty"`
}

func (err Error) Error() string {
	var locationStr string
	if err.File != "" || err.Function != "" {
		locationStr = fmt.Sprintf("%s:%d %s", err.File, err.Line, err.Function)
	}
	return fmt.Sprintf("<%d> %s: %v (%s)", err.ID, err.Description, err.Cause, locationStr)
}

func (err *Error) MarshalJSON() ([]byte, error) {
	type ShadowError Error
	type Default struct {
		*ShadowError
		CauseMsg string `json:"cause,omitempty"`
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
		if err.Cause != nil {
			return json.Marshal(&Default{
				ShadowError: (*ShadowError)(err),
				CauseMsg:    err.Cause.Error(),
			})
		} else {
			return json.Marshal(&Default{
				ShadowError: (*ShadowError)(err),
			})
		}
	}
}

var errorDescriptions = map[ErrorID]string{}

func NewError(id ErrorID, cause error) *Error {
	return newError(id, cause, false)
}

func NewErrorSkipFunc(id ErrorID, cause error) *Error {
	return newError(id, cause, true)
}

func newError(id ErrorID, cause error, skipCaller bool) *Error {
	// Skip only this func or one more for location
	skipPC := 1
	if skipCaller {
		skipPC++
	}

	pc, file, line, _ := runtime.Caller(skipPC)
	var functionName string
	if function := runtime.FuncForPC(pc); function != nil {
		functionName = function.Name()
	}
	return &Error{
		ID:          id,
		Description: errorDescriptions[id],
		Cause:       cause,
		File:        file,
		Line:        line,
		Function:    functionName,
	}
}
