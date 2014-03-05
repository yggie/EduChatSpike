package xmpp

type Connection interface {
  Send(message string) error
  Push(message string)
  UserName() string
  JID() string
}
