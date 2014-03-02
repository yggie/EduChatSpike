package sasl

import(
  "fmt"
  "log"
  "errors"
  "strings"
  "github.com/yggie/EduChatSpike/modules/auth"
)

type metadata struct {
  Mechanism string
  Data []byte
}

var (
  AuthFailError = errors.New("Authentication failed\n")

  store = make(map[string]metadata)
)

/// TODO guarantee uniqueness across application
func genNonce() []byte {
  return []byte("MsQUY9iw0T9fx2MUEz6LZPwGuhVvWAhc")
}

func Init(mechanism string, sID string, data []byte) ([]byte, error) {
  switch mechanism {
  case "SCRAM-SHA-1":
    msg, err := auth.DecodeBase64(data)
    if err != nil {
      log.Println(err)
      return nil, AuthFailError
    }
    nonce := genNonce()
    response, err := InitialResponseSCRAMSHA1(msg, nonce)
    if err != nil {
      return nil, err
    }
    splitMsg := strings.Split(string(msg), ",")
    store[sID] = metadata{
      Mechanism: "SCRAM-SHA-1",
      Data: []byte(splitMsg[2] + "," + splitMsg[3] + "," + string(response)),
    }
    return auth.EncodeBase64(response), nil

  default:
    return nil, fmt.Errorf("unsupported authentication mechanism: %s\n", mechanism)
  }

  log.Printf("Impossible program flow!\n")
  return nil, nil
}

func Respond(sID string, data []byte) ([]byte, error) {
  meta, ok := store[sID]
  if !ok {
    return nil, AuthFailError
  }

  switch meta.Mechanism {
  case "SCRAM-SHA-1":
    raw, err := auth.DecodeBase64(data)
    if err != nil {
      log.Println(err)
      return nil, AuthFailError
    }

    response, err := FinalResponseSCRAMSHA1(raw, meta.Data)
    if err != nil {
      return nil, err
    }

    return auth.EncodeBase64(response), nil

  default:
    return nil, fmt.Errorf("unsupported authentication mechanism: %s\n", meta.Mechanism)
  }
}

func WasVerified(sID string) bool {
  _, ok := store[sID]
  if ok {
    delete(store, sID)
    return true
  }
  return false
}


