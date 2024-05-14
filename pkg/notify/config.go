package notify

import (
	"dussh/pkg/notify/notification"
	"dussh/pkg/notify/provider"
	"dussh/pkg/notify/provider/email"
)

type Config struct {
	// Email is the configuration for the email notify provider
	Email *email.NotificationProvider
}

func (c *Config) GetNotificationProviderByType(
	notificationType notification.Type,
) provider.NotificationProvider {
	switch notificationType {
	case notification.TypeEmail:
		return c.Email
	default:
		return nil
	}
}
