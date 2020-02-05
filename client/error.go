package client

import (
	"fmt"
)

// Error ...
type Error struct {
	Code    int    `json:"code,string,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
