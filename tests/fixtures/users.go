package fixtures

import (
  "github.com/yggie/EduChatSpike/modules/auth"
  "github.com/yggie/EduChatSpike/modules/models"
  "github.com/yggie/EduChatSpike/modules/records"
)

var (
  StubUser = models.NewUser("user", "secret")
)

type StubbedUserFinder struct {
  records.UserFinder
}

func (s StubbedUserFinder) FindByName(name string) models.User {
  return StubUser
}

func CreateStubUserWithSalt(name, password string, salt []byte) *models.User {
  StubUser = models.User{
    Name: name,
    Pass: auth.Pass{
      Salt: salt,
      Key: auth.Hi([]byte(password), salt, auth.PASSWORD_ITERATIONS),
    },
  }
  return &StubUser
}

