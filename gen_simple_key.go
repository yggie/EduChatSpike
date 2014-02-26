package main

import (
  "io"
  "fmt"
  "crypto/sha1"
  "yggie/EduChatSpike/lib/auth"
)

func main() {
  sha := sha1.New()
  io.WriteString(sha, "AgileVentures EduChat")
  salt := sha.Sum(nil)
  fmt.Printf("%x\n", string(salt))

  pwd := "embeddedchatforall"
  key := auth.Hi([]byte(pwd), salt, 8080)
  fmt.Printf("%x\n", key)
}
