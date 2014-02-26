package auth

import (
  "strconv"
  "strings"
  "crypto/sha1"
  "crypto/hmac"
  "encoding/base64"
  "yggie/EduChatSpike/lib/models"
)

const (
  SCRAM_SHA_1_Iterations = 8080
)

/// The Hi function defined in RFC5802
func Hi(pass []byte, salt []byte, iters int) []byte {
  u := append(salt, 0x31, 0x00, 0x00, 0x00)
  res := make([]byte, sha1.Size)
  for i := 0; i < iters; i++ {
    h := hmac.New(sha1.New, u)
    u = h.Sum(nil)
    for j := 0; j < sha1.Size; j++ {
      res[j] ^= u[j]
    }
  }

  return res
}

/// Responds to the Client's first message as defined in SASL SCRAM-SHA-1, RFC5802
func scramFirstMessage(clientFirstMessage []byte) ([]byte, error) {
  decoded := make([]byte, 255)
  _, err := base64.StdEncoding.Decode(decoded, clientFirstMessage)

  if err != nil {
    return nil, AuthFailError
  }

  items := strings.Split(string(decoded), ",")

  // expect exactly 4 components
  if len(items) != 4 {
    return nil, AuthFailError
  }

  var salt, r string

  for _, item := range items {
    s := strings.Split(item, "=")
    // we are expecting a format of [var]=[value]
    if len(s) == 2 {
      switch s[0] {
      case "n":
        user := models.GetUser(s[1])
        salt = string(user.Pass.Salt)

      case "r":
        nonce, err := newNonce()
        if err != nil {
          return nil, err
        }
        r = s[1] + string(nonce)

      default:
        return nil, AuthFailError
      }
    }
  }

  if len(salt) == 0 || len(r) == 0 {
    return nil, AuthFailError
  }

  // response should be: r=[ClientNonce + ServerNonce],s=[UserPassSalt],i=[Iterations]
  response := "r=" + r + ",s=" + salt + ",i=" + strconv.Itoa(SCRAM_SHA_1_Iterations)

  result := make([]byte, 255)
  base64.StdEncoding.Encode(result, []byte(response))

  return result, nil
}

