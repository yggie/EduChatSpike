package security_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/yggie/EduChatSpike/modules/records"
  "github.com/yggie/EduChatSpike/modules/models"

  "testing"
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

func TestAuth(t *testing.T) {
  records.Users = StubbedUserFinder{}
  RegisterFailHandler(Fail)
  RunSpecs(t, "Security Features Suite")
}

