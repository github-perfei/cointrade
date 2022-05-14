package push

type pusher interface {
	SendMessage(message string) error
}