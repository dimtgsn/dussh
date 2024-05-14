package notification

type Notification struct {
	Type        Type
	ContentType ContentType
	To          []string
	Subject     string
	Body        string
}
