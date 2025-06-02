package zhttp

import (
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewStatusCodeError(statusCode int) *StatusCodeError {
	return &StatusCodeError{statusCode: statusCode}
}

type StatusCodeError struct {
	statusCode int // e.g. 200
}

func (s *StatusCodeError) Error() string {
	return fmt.Sprintf("http status code %d", s.statusCode)
}

func (s *StatusCodeError) StatusCode() int {
	return s.statusCode
}

func (s *StatusCodeError) Is(target error) bool {
	//nolint:errorlint
	if other, is := target.(*StatusCodeError); is {
		return s.statusCode == other.statusCode
	}

	return false
}
