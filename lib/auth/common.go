package auth

import (
  "fmt"
  "log"
  "errors"
  "crypto/rand"
)

var (
  AuthFailError = errors.New("Authentication failed\n")
)

type Pass struct {
  Salt []byte
  Key []byte
}

/// TODO guarantee uniqueness across application
func newNonce() ([]byte, error) {
  nonce := make([]byte, 50)
  _, err := rand.Read(nonce)
  return nonce, err
}

func Init(mechanism string, sID string, data []byte) ([]byte, error) {
  switch mechanism {
  case "SCRAM-SHA-1":
    return scramFirstMessage(data)

  default:
    // TODO write error result
    return nil, fmt.Errorf("unsupported authentication mechanism: %s\n", mechanism)
  }

  log.Printf("Impossible program flow!\n")
  return nil, nil
}

func Respond(sID string, data []byte) ([]byte, error) {
  return nil, nil
}
