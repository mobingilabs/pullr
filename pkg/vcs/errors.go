package vcs

import "errors"

// Errors constants
var (
	ErrInvalidWebhook        = errors.New("Invalid webhook request")
	ErrInvalidWebhookPayload = errors.New("Invalid webhook payload")
	ErrUnsupportedEvent      = errors.New("Unsupported event")
)
