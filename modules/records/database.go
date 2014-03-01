package records

import (
  "github.com/yggie/EduChatSpike/modules/models"
)

type UserFinder interface {
  FindByName(name string) models.User
}

type Database struct {
  Users UserFinder
}

