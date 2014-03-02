package sasl

import (
  "log"
  "strconv"
  "strings"
  "github.com/yggie/EduChatSpike/modules/auth"
  "github.com/yggie/EduChatSpike/modules/records"
)

func InitialResponseSCRAMSHA1(initialMsg, nonce []byte) ([]byte, error) {
  // first message must start with either "n", "y" or "p"
  if initialMsg[0] != 'n' { // && initialMsg[0] != 'p' && initialMsg[0] != 'y' {
    log.Println("invalid initial response")
    return nil, AuthFailError
  }

  s := strings.Split(string(initialMsg), ",")

  // well formed messages will result in 4 components
  if len(s) != 4 {
    log.Println("invalid initial response")
    return nil, AuthFailError
  }

  username := string(s[2][2:])
  clientNonce := string(s[3][2:])

  user := records.Users.FindByName(username)

  encodedSalt := string(auth.EncodeBase64(user.Pass.Salt))
  return []byte("r=" + clientNonce + string(nonce) + ",s=" + encodedSalt + ",i=" + strconv.Itoa(auth.PASSWORD_ITERATIONS)), nil
}

func FinalResponseSCRAMSHA1(clientFinalMsg, prevMsg []byte) ([]byte, error) {
  oldMsg := strings.Split(string(prevMsg), ",")
  components := strings.Split(string(clientFinalMsg), ",")

  if oldMsg[2] != components[1] {
    log.Println("session nonce does not match")
    return nil, AuthFailError
  }

  authMsg := string(prevMsg) + "," + components[0] + "," + components[1]

  username := oldMsg[2][2:]
  user := records.Users.FindByName(username)

  // for some reason there are null characters at the end?
  tmp := strings.TrimRight(components[2][2:], "\x00")
  clientProof, err := auth.DecodeBase64([]byte(tmp))
  if err != nil {
    log.Println(err)
    return nil, AuthFailError
  }
  clientKey := auth.HMAC(user.Pass.Key, "Client Key")
  storedKey := auth.H(clientKey)
  clientSignature := auth.HMAC(storedKey, authMsg)

  length := len(storedKey)
  authFailed := false
  for i := 0; i < length; i++ {
    v := clientKey[i] ^ clientSignature[i]
    if clientProof[i] != v {
      authFailed = true
    }
  }

  // avoid early exit when comparing password related data
  if authFailed {
    log.Println("Invalid Client Proof")
    return nil, AuthFailError
  }

  serverKey := auth.HMAC(user.Pass.Key, "Server Key")
  signature := auth.HMAC(serverKey, authMsg)
  armouredSignature := auth.EncodeBase64(signature)
  return []byte("v=" + string(armouredSignature)), nil
}

