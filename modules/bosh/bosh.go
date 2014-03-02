package bosh

import (
  "fmt"
  "log"
  "net/http"
)

var (
  InvalidSessionError = fmt.Errorf("invalid session id")
)

func HttpBindHandler(w http.ResponseWriter, r *http.Request) {
  // parse request body
  request, err := ParseRequest(r)

  if err == nil {

    listener := Listener{
      Request: request,
      Writer: w,
    }

    err = request.Process(&listener)
  }

  if err != nil {
    log.Println(err)
    w.WriteHeader(400)
  }
}

