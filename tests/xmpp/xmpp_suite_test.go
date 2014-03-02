package xmpp_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/yggie/EduChatSpike/modules/xmpp"

  "testing"
)

func TestXMPP(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "XMPP Processing Suite")
}

type Result struct {
  Stanza *Stanza
  Error error
  Response string
  Processed bool
}

func (r *Result) Send(str string) {
  var err error
  r.Stanza, err = ParseString(str)
  if err != nil {
    panic("parsing error in xml provided by the test, something must be wrong")
  }
  r.Error = r.Stanza.Process(r)
}

func (r *Result) InitSASL(auth *Auth) error {
  panic("STUB")
  return nil
}

func (r *Result) RespondSASL(clientResponse *Response) error {
  panic("STUB")
  return nil
}

func (r *Result) Unprocessed() error {
  r.Processed = false
  return nil
}

func (r *Result) Complete(response string) error {
  r.Response = response
  r.Processed = true
  return nil
}

