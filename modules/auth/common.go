package auth

import (
  "crypto/hmac"
  "crypto/sha1"
  "encoding/base64"
)

const (
  PASSWORD_ITERATIONS = 4096
)

type Pass struct {
  Salt []byte
  Key  []byte
}

func GenPass(password string) Pass {
  salt := []byte("b977d0396381a0d2092d04075b54c91e90442d06")
  return Pass{
    Salt: salt,
    Key: Hi([]byte(password), salt, PASSWORD_ITERATIONS),
  }
}

func EncodeBase64(data []byte) []byte {
  dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
  base64.StdEncoding.Encode(dst, data)
  return dst
}

func DecodeBase64(raw []byte) ([]byte, error) {
  dst := make([]byte, base64.StdEncoding.DecodedLen(len(raw)))
  _, err := base64.StdEncoding.Decode(dst, raw)
  if err != nil {
    return nil, err
  }
  return dst, nil
}

// The H function defined in RFC5802
func H(data []byte) []byte {
  h := sha1.New()
  h.Write(data)
  return h.Sum(nil)
}

// The HMAC function defined in RFC5802
func HMAC(key []byte, str string) []byte {
  h := hmac.New(sha1.New, key)
  h.Write([]byte(str))
  return h.Sum(nil)
}

/// The Hi function defined in RFC5802
func Hi(pass []byte, salt []byte, iters int) []byte {
  u := append(salt, 0x00, 0x00, 0x00, 0x01)
  res := make([]byte, sha1.Size)
  for i := 0; i < iters; i++ {
    u = HMAC(pass, string(u))
    for j := 0; j < sha1.Size; j++ {
      res[j] ^= u[j]
    }
  }

  return res
}

