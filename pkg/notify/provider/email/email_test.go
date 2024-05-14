package email

import (
	"context"
	"dussh/pkg/notify/notification"
	"testing"
)

func TestSend(t *testing.T) {
	var (
		provider = &NotificationProvider{
			From:      "example@conmpany.com",
			Username:  "f0d73492594ef7",
			Password:  "ee6d5d7a04bbbd",
			Host:      "sandbox.smtp.mailtrap.io",
			Port:      2525,
			TLSEnable: false,
		}
		ntf = &notification.Notification{
			Type:        notification.TypeEmail,
			ContentType: notification.ContentTypePlain,
			To:          []string{"gasanyandmitry@yandex.ru"},
			Subject:     "Test",
			Body:        "Test",
		}
	)

	if !provider.IsValid() {
		t.Fatal("provider config is invalid")
	}
	if err := provider.Send(context.Background(), ntf); err != nil {
		t.Fatalf("failed to send notify: %v", err)
	}

	t.Log("send notify successfully")
}
