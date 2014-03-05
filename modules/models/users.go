package models

import (
  "github.com/yggie/EduChatSpike/modules/auth"
)

type User struct {
  Name string
  Pass auth.Pass
}

func NewUser(username, password string) User {
  return User{
    Pass: auth.GenPass(password),
  }
}

