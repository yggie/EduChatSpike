package fixtures

import (
  "github.com/yggie/EduChatSpike/modules/xmpp"
)

type ConnectionFixture struct {
  Connection xmpp.Connection
  Messages []string
}

func (f *ConnectionFixture) Send(message string) error {
  return f.Connection.Send(message)
}

func (f *ConnectionFixture) Receive(message string) {
  f.Messages = append(f.Messages, message)
}

func (f *ConnectionFixture) Bind(conn xmpp.Connection) {
  f.Connection = conn
}

func (f *ConnectionFixture) ClearMessages() {
  f.Messages = nil
}

func (f *ConnectionFixture) LastMessage() string {
  length := len(f.Messages)
  if length == 0 {
    return ""
  } else {
    return f.Messages[length-1]
  }
}

func NewAnonymousConnection() *ConnectionFixture {
  fixture := ConnectionFixture{}
  fixture.Connection = &xmpp.Anonymous{
    Receiver: &fixture,
    Binder: &fixture,
  }
  return &fixture
}
