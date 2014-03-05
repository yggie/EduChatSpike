package xmpp_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/yggie/EduChatSpike/tests/fixtures"
)

var _ = Describe("Info queries <iq/>", func() {

  Context("response to pings", func() {

    var (
      fixture *ResourceFixture
      err error
    )

    BeforeEach(func() {
      fixture = LogInAsUser("jack")
      err = fixture.Send("<iq to='educhat.spike' id='ping1' xmlns='jabber:client'><ping xmlns='urn:xmpp:ping'/></iq>")
    })

    It("should not have any errors", func() {
      Expect(err).To(BeNil())
    })

    It("should respond with a type of 'result'", func() {
      Expect(fixture.LastMessage()).To(MatchRegexp("type='result'"))
    })

    It("should return the result with the same id used in the query", func() {
      Expect(fixture.LastMessage()).To(MatchRegexp("id='ping1'"))
    })

    It("should be correctly addressed to the client", func() {
      Expect(fixture.LastMessage()).To(MatchRegexp("to='" + fixture.JID() + "'"))
    })
  })
})
