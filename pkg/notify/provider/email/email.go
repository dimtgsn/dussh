package email

import (
	"context"
	"crypto/tls"
	"dussh/pkg/notify/notification"
	"gopkg.in/gomail.v2"
	"math"
	"strings"
)

type NotificationProvider struct {
	From      string
	Username  string
	Password  string
	Host      string
	Port      int
	TLSEnable bool
}

// IsValid returns whether the provider's configuration is valid
func (provider *NotificationProvider) IsValid() bool {
	isValid := len(provider.From) > 0 && len(provider.Host) > 0 &&
		provider.Port > 0 && provider.Port < math.MaxUint16

	return isValid
}

// Send a notification using the provider
func (provider *NotificationProvider) Send(
	ctx context.Context,
	notification *notification.Notification,
) error {
	var username string
	if len(provider.Username) > 0 {
		username = provider.Username
	} else {
		username = provider.From
	}
	m := gomail.NewMessage()
	m.SetHeader("From", provider.From)
	m.SetHeader("To", strings.Join(notification.To, ","))
	m.SetHeader("Subject", notification.Subject)
	m.SetBody(string(notification.ContentType), notification.Body)
	var d *gomail.Dialer
	if len(provider.Password) == 0 {
		// Get the domain in the From address
		localName := "localhost"
		fromParts := strings.Split(provider.From, `@`)
		if len(fromParts) == 2 {
			localName = fromParts[1]
		}
		// Create a dialer with no authentication
		d = &gomail.Dialer{Host: provider.Host, Port: provider.Port, LocalName: localName}
	} else {
		// Create an authenticated dialer
		d = gomail.NewDialer(provider.Host, provider.Port, username, provider.Password)
	}
	if !provider.TLSEnable {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return d.DialAndSend(m)
}
