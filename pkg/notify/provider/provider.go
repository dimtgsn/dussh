package provider

import (
	"context"
	"dussh/pkg/notify/notification"
	"dussh/pkg/notify/provider/email"
)

type NotificationProvider interface {
	// IsValid returns whether the provider's configuration is valid
	IsValid() bool

	// Send a notification using the provider
	Send(context.Context, *notification.Notification) error
}

var (
	_ NotificationProvider = (*email.NotificationProvider)(nil)
)
