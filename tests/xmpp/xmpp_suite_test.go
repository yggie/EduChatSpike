package xmpp_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/yggie/EduChatSpike/tests/fixtures"

  "testing"
)

func TestXMPP(t *testing.T) {
  fixtures.Initialize()
  RegisterFailHandler(Fail)
  RunSpecs(t, "XMPP Processing Suite")
}

