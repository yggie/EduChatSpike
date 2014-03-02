package xmpp_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Info Queries", func() {

  Context("ping", func() {
    It("should respond to a ping", func() {
      r := Result{}
      r.Send("<iq to='examples.org' type='get' id='ping1' xmlns='jabber:client'><ping xmlns='urn:xmpp:ping' /></iq>")
      Expect(r.Response).To(Equal("")) // this is wrong
    })
  })
})
