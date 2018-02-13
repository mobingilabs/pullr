package vcs

import "errors"

// Errors constants
var (
	ErrAuthRequired          = errors.New("vcs client needs to authenticate")
	ErrInvalidWebhook        = errors.New("invalid webhook request")
	ErrInvalidWebhookPayload = errors.New("invalid webhook payload")
	ErrUnsupportedEvent      = errors.New("unsupported event")
)
