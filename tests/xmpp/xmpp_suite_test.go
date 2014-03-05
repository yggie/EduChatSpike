package xmpp_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestXMPP(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "XMPP Processing Suite")
}

