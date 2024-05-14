package notification

type Type string

const (
	// TypeEmail is Type for the email notify provider
	TypeEmail Type = "email"
)

type ContentType string

const (
	ContentTypePlain ContentType = "text/plain"
	ContentTypeHTML  ContentType = "text/html"
)
