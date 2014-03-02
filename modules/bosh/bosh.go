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

    processor := RequestProcessor{
      Request: request,
      Writer: w,
    }

    err = request.Process(&processor)
  }

  if err != nil {
    log.Println(err)
    w.WriteHeader(400)
  }
}

