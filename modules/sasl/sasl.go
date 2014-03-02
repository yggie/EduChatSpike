package sasl

import(
  "log"
  "errors"
  "strings"
  "github.com/yggie/EduChatSpike/modules/auth"
)

type UnsupportedAuthenticationMechanism string
func (mech UnsupportedAuthenticationMechanism) Error() string {
  return "unsupported authentication mechanism: [" + string(mech) + "]"
}

var (
  AuthFailError = errors.New("Authentication failed\n")
)

/// TODO guarantee uniqueness across application
func genNonce() []byte {
  return []byte("MsQUY9iw0T9fx2MUEz6LZPwGuhVvWAhc")
}

// Initializes the SASL handshake using the specified mechanism
func Init(mechanism string, data []byte) ([]byte, []byte, error) {
  switch mechanism {
  case "SCRAM-SHA-1":
    msg, err := auth.DecodeBase64(data)
    if err != nil {
      log.Println(err)
      return nil, nil, AuthFailError
    }
    nonce := genNonce()
    response, err := InitialResponseSCRAMSHA1(msg, nonce)
    if err != nil {
      return nil, nil, err
    }
    splitMsg := strings.Split(string(msg), ",")
    data := []byte(splitMsg[2] + "," + splitMsg[3] + "," + string(response))

    return auth.EncodeBase64(response), data, nil

  default:
    return nil, nil, UnsupportedAuthenticationMechanism(mechanism)
  }

  log.Printf("Impossible program flow!\n")
  return nil, nil, nil
}

// Continues the SASL handshake
func Respond(mech string, data []byte, storedData []byte) ([]byte, error) {
  switch mech {
  case "SCRAM-SHA-1":
    raw, err := auth.DecodeBase64(data)
    if err != nil {
      log.Println(err)
      return nil, AuthFailError
    }

    response, err := FinalResponseSCRAMSHA1(raw, storedData)
    if err != nil {
      return nil, err
    }

    return auth.EncodeBase64(response), nil

  default:
    return nil, UnsupportedAuthenticationMechanism(mech)
  }
}
