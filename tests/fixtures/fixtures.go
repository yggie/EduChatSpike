package fixtures

import (
  . "github.com/yggie/EduChatSpike/modules/xmpp"
  "github.com/yggie/EduChatSpike/modules/records"
)

func Initialize() {
  records.Users = StubbedUserFinder{}
}

type ResourceFixture struct {
  UserResource *UserResource
  Messages []string
}

func (r *ResourceFixture) Send(message string) error {
  return r.UserResource.Send(message)
}

func (r *ResourceFixture) Receive(message string) {
  r.Messages = append(r.Messages, message)
}

func (r *ResourceFixture) LastMessage() string {
  length := len(r.Messages)
  if length == 0 {
    return ""
  } else {
    return r.Messages[length-1]
  }
}

func (r *ResourceFixture) JID() string {
  return r.UserResource.JID()
}

func LogInAsUser(name string) *ResourceFixture {
  fixture := ResourceFixture{}
  fixture.UserResource = NewResourceForUser(name, &fixture)
  return &fixture
}

