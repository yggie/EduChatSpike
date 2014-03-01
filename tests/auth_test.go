package security_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  . "github.com/yggie/EduChatSpike/modules/auth"
  "encoding/hex"
)


var _ = Describe("hashing functions", func() {

  It("should correctly perform the RSC-5802 H function", func() {
    raw := H([]byte("water"))
    hexResult := "6d5a45920a15adea049c8f22d569ff209625a43b"
    hex, _ := hex.DecodeString(hexResult)
    Ω(raw).To(Equal(hex))
  })

  It("should correctly perform the RSC-5802 HMAC function", func() {
    raw := HMAC([]byte("the key"), "the message")
    Ω(hex.EncodeToString(raw)).To(Equal("9f441d3ad40a069c1a82e21007f3e52fd6195b62"))
  })

  It("should correctly perform the RSC-5802 Hi function", func() {
    raw := Hi([]byte("secret key"), []byte("salty"), 4096)
    Expect(hex.EncodeToString(raw)).To(Equal("dfe1c5ec99f2246d7b147b808622390679fe14cb"))
  })
})

