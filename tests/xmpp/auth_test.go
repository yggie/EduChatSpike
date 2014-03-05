package xmpp_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/yggie/EduChatSpike/tests/fixtures"

  "github.com/yggie/EduChatSpike/modules/auth"
)

var _ = Describe("Anonymous connections", func() {

  Context("authentication", func() {

    salt, _  := auth.DecodeBase64([]byte("QSXCR+Q6sek8bf92"))
    user := CreateStubUserWithSalt("user", "pencil", salt)
    connection := NewAnonymousConnection()

    BeforeEach(func() {
      connection.ClearMessages()
    })

    It("should respond to the first message in the authentication handshake", func() {
      clientMsg := "n,,u=" + user.Name + ",n=fyko+d2lbbFgONRv9qkxdawL"
      msg := auth.EncodeBase64([]byte(clientMsg))
      err := connection.Send("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='SCRAM-SHA-1'>" + string(msg) + "</auth>")
      if err != nil {
        panic(err)
      }

      Expect(connection.LastMessage()).To(MatchRegexp("^<challenge(.*)challenge>$"))
    })

    It("should respond with successful authentication on the second part of the handshake", func() {
      clientFinalMsg := "c=biws,r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,p=v0X8v3Bz2T0CJGbJQyF0X+HI4Ts="
      msg := auth.EncodeBase64([]byte(clientFinalMsg))
      err := connection.Send("<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" + string(msg) + "</response>")
      if err != nil {
        panic(err)
      }

      Expect(connection.LastMessage()).To(MatchRegexp("^<success(.*)sucess>$"))
    })
  })
})
