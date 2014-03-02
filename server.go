package main

import (
  "log"
  "net/http"
  "html/template"
  "github.com/yggie/EduChatSpike/modules/bosh"
  "github.com/yggie/EduChatSpike/modules/models"
  "github.com/yggie/EduChatSpike/modules/records"
)

func main() {
  // setup user finder
  records.Users = StubbedUserFinder{}

  http.HandleFunc("/", handler)
  http.HandleFunc("/http-bind", bosh.HttpBindHandler)
  http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("tmp/public/"))))
  log.Printf("Listening in on port 3000\n")
  err := http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatal(err)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("tmp/views/index.html")
  log.Printf("Received request from %s\n", r.URL.Host + r.URL.Path)
  t.Execute(w, nil)
}

type StubbedUserFinder struct {
  records.UserFinder
}

func (s StubbedUserFinder) FindByName(name string) models.User {
  return models.NewUser("educhatspikeuser", "embeddedchatforall")
}
