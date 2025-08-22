package domain

import "errors"

type SMS struct {
	To      string
	Message string
}

var (
	ErrCircuitOpen            = errors.New("circuit breaker is open")
	ErrMaxRetries             = errors.New("max retries exceeded")
	ErrServiceTimeout         = errors.New("sms service timeout")
	ErrExternalServiceFailure = errors.New("external service failure")
)
