package bosh

import (
  "fmt"
  "log"
  "time"
  "net/http"
)

var (
  InvalidSessionError = fmt.Errorf("invalid session id")
)

func HttpBindHandler(w http.ResponseWriter, r *http.Request) {
  // parse request body
  request, err := ParseRequest(r)
  if err != nil {
    log.Println(err)
    w.WriteHeader(400)
    return
  }

  // set up default reply headers
  w.Header().Set("Content-Type", "text/xml")

  response, err := respondTo(request)
  if err != nil {
    log.Println(err)
    return
  }

  fmt.Printf("\nResponse:\n%s\n\n", response)
  w.Write([]byte(response))
}

func respondTo(r *Request) (string, error) {
  if r.SID == "" {
    return createNewSession(r), nil

  } else if r.HasInfoQuery() {
    return parseData(r)

  } else if r.HasAuth() {
    return initSASL(r)

  } else if r.HasResponse() {
    return respondSASL(r)

  } else if r.ShouldRestart() {
    return restartSession(r)

  } else {
    time.Sleep(60 * time.Second)
  }

  return "", nil
}

