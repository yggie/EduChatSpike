package main

import (
  "log"
  "./lib/conn"
  "net/http"
  "html/template"
)

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/http-bind", conn.HttpBindHandler)
  http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
  log.Printf("Listening in on port 3000\n")
  err := http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatal(err)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("views/index.html")
  log.Printf("Received request from %s\n", r.URL.Host + r.URL.Path)
  t.Execute(w, nil)
}
