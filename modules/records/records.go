package records

import (
  "github.com/yggie/EduChatSpike/modules/models"
)

var (
  Users UserFinder
)

type UserFinder interface {
  FindByName(name string) models.User
}

