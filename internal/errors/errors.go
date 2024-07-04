package errors

import (
	"fmt"
)

type ConfigError struct {
	Message string
}

func (e ConfigError) Error() string {
	return e.Message
}

func NewConfigError(format string, a ...interface{}) *ConfigError {
	return &ConfigError{
		Message: fmt.Sprintf(format, a...),
	}
}

type RequestError struct {
	Message string
}

func (e RequestError) Error() string {
	return e.Message
}

func NewRequestError(format string, a ...interface{}) *RequestError {
	return &RequestError{
		Message: fmt.Sprintf(format, a...),
	}
}
