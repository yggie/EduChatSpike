package xmpp

import (
  "encoding/xml"
)

// represents a single XMPP packet that could be received
type Stanza struct {
  // if exists, it is the client's preferred authentication mechanism
  Auth []Auth           `xml:"auth"`

  // if exists, it is the client response to an SASL challenge
  Response []Response   `xml:"response"`

  // if exists, represents info queries from the client
  Iq []Iq               `xml:"iq"`
}

func (s *Stanza) HasAuth() bool {
  return len(s.Auth) != 0
}

func (s *Stanza) HasResponse() bool {
  return len(s.Response) != 0
}

func (s *Stanza) HasInfoQuery() bool {
  return len(s.Iq) != 0
}

func ParseString(str string) (*Stanza, error) {
  stanza := Stanza{}
  err := xml.Unmarshal([]byte(str), &stanza);
  if err != nil {
    return nil, err
  }

  return &stanza, nil
}
