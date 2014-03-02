package security_test

import (
  "strconv"
  "strings"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/yggie/EduChatSpike/modules/sasl"
  "github.com/yggie/EduChatSpike/modules/auth"
  "github.com/yggie/EduChatSpike/modules/models"
)

var _ = Describe("SCRAM-SHA-1", func() {

  Describe("first message", func() {

    var (
      clientNonce string
      username string
      clientFirstMsg []byte
      nonce []byte
    )

    BeforeEach(func() {
      clientNonce = "MsQUY9iw0T9fx2MUEz6LZPwGuhVvWAhc"
      username = "chris"
      nonce = []byte("MsQUY9iw0T9fx2MUEz6LZPwGuhVvWAhc")

      clientFirstMsg = []byte("n,,n=" + username + ",r=" + clientNonce)
    })

    Context("responding to a valid first message", func() {

      var (
        msg string
        err error
        components []string
      )

      BeforeEach(func() {
        var raw []byte
        raw, err = InitialResponseSCRAMSHA1(clientFirstMsg, nonce)
        msg = string(raw)
        components = strings.Split(msg, ",")
      })

      It("should not generate an error", func() {
        Ω(err).To(BeNil())
      })

      It("should have exactly 3 components", func() {
        Ω(len(components)).To(Equal(3))
      })

      It("should set the first component to be the combined nonce", func() {
        Ω(components[0]).To(Equal("r=" + string(clientNonce) + string(nonce)))
      })

      It("should set the second component to the users salt", func() {
        raw := StubUser.Pass.Salt
        encBase64 := auth.EncodeBase64(raw)
        Ω(components[1][2:]).To(Equal(string(encBase64)))
      })

      It("should contain the correct number of iterations", func() {
        iters, _ := strconv.Atoi(components[2][2:])
        Ω(iters).To(Equal(auth.PASSWORD_ITERATIONS))
      })
    })

    Context("with invalid requests", func() {
      It("should throw an error when it receives malformed requests", func() {
        _, err := InitialResponseSCRAMSHA1([]byte("thisissomerandommessage"), nonce)
        Expect(err).NotTo(BeNil())
      })
    })
  })

  Describe("final response", func() {
    It("should not accept mismatch between the input and cached nonce", func() {
      _, err := FinalResponseSCRAMSHA1([]byte("c=,r=abcd,p="), []byte("n=,r=,r=abc,s=,i="))
      Expect(err).NotTo(BeNil())
    })

    It("should not accept invalid Client Proofs", func() {
      _, err := FinalResponseSCRAMSHA1([]byte("c=,r=abc,p=hello"), []byte("n=,r=,r=abc,s=,i="))
      Expect(err).NotTo(BeNil())
    })
  })

  Describe("full handshake", func() {

    salt, _  := auth.DecodeBase64([]byte("QSXCR+Q6sek8bf92"))
    StubUser = models.User{
      Pass: auth.Pass{
        Salt: salt,
        Key: auth.Hi([]byte("pencil"), salt, auth.PASSWORD_ITERATIONS),
      },
    }
    iters := strconv.Itoa(auth.PASSWORD_ITERATIONS)
    clientNonce := "fyko+d2lbbFgONRv9qkxdawL"
    nonce := "3rfcNHYJY1ZVvWVs7j"
    username := "user"

    It("should respond correctly to the first message", func() {
      clientMsg := "n,,u=" + username + ",n=" + clientNonce
      msg, err := InitialResponseSCRAMSHA1([]byte(clientMsg), []byte(nonce))

      Expect(err).To(BeNil())
      Expect(string(msg)).To(Equal("r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,s=QSXCR+Q6sek8bf92,i=" + iters))
    })

    It("should respond correctly to the final message", func() {
      prevMsg := "n=" + username + ",r=" + clientNonce + ",r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,s=QSXCR+Q6sek8bf92,i=" + iters
      // should also be able to handle cases with null characters
      clientFinalMsg := "c=biws,r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,p=v0X8v3Bz2T0CJGbJQyF0X+HI4Ts=\x00"
      msg, _ := FinalResponseSCRAMSHA1([]byte(clientFinalMsg), []byte(prevMsg))

      // Expect(err).To(BeNil())
      Expect(string(msg)).To(Equal("v=rmF9pqV8S7suAoZWja4dJRkFsKQ="))
    })
  })
})
