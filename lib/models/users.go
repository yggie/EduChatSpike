package models

type Pass struct {
  Salt []byte
  Key  []byte
}

type User struct {
  Pass Pass
}

func GetUser(string) *User {
  return &User{
    Pass: Pass{
      Salt: []byte("b977d0396381a0d2092d04075b54c91e90442d06"),
      Key: []byte("e147154acfcc582bf3989b49266ff094c45498f2"),
    },
  }
}
