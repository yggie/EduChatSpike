package xmpp

import (
  "strings"
  "github.com/yggie/EduChatSpike/modules/auth"
)

type Auth struct {
  Namespace string      `xml:"xmlns,attr"`
  Mechanism string      `xml:"mechanism,attr"`
  Data []byte           `xml:",chardata"`
  XMLName interface{}   `xml:"auth"`
}

func (a *Auth) Exists() bool {
  return len(a.Namespace) != 0
}

func (a *Auth) GetUserName() (string, error) {
  raw, err := auth.DecodeBase64(a.Data)
  if err != nil {
    return "", err
  }

  parts := strings.Split(string(raw), ",")

  return parts[2][2:], nil
}
