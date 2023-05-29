// connection package handle connection to rest-api server Warships Online API
package connection

import "fmt"

// ErrorMessage Store server message when response status code is 403 usualy "session not found"
type ErrorMessage struct {
	Message string `json:"message"`
}

// ToekenError error for wrong token
type TokenError struct {
	Token string
}

func (e *TokenError) Error() string {
	return fmt.Sprintf("invalid token: %s", e.Token)
}

// RequestError error for status code not ok
type RequestError struct {
	StatusCode int
	Err        string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", e.StatusCode, e.Err)
}
