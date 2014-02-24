package main

import (
  "fmt"
  "log"
  "time"
  "net/http"
  "io/ioutil"
  "html/template"
  "encoding/xml"
)

/// The session creation request body as defined by XEP-0206
type Request struct {
  Content string        `xml:"content,attr"`
  From string           `xml:"from,attr"`
  Hold  int             `xml:"hold,attr"`
  RID int64             `xml:"rid,attr"`
  To  string            `xml:"to,attr"`
  Route string          `xml:"route,attr"`
  Secure bool           `xml:"secure,attr"`
  Wait  int             `xml:"wait,attr"`
  Lang  string          `xml:"xml:lang,attr"`
  XMPPVersion string    `xml:"xmpp:version,attr"`
  Namespace string      `xml:"xmlns,attr"`
  XMPPNamespace  string `xml:"xmlns:xmpp,attr"` 
  // Strophe version??
  // Version string      `xml:"ver,attr"`
}

type Mechanisms struct {
  Mechanisms []string      `xml:"mechanism"`
}

type Bind struct {
  Namespace string        `xml:"xmlns,attr"`
}

type Feature struct {
  // Mechanisms Mechanisms   `xml:"mechanisms"`
  Bind Bind               `xml:"bind"`
}

type Success struct {
  Namespace string        `xml:"xmlns,attr"`
}

type Body struct {
  Wait int                `xml:"wait,attr"`
  Inactivity int          `xml:"inactivity,attr"`
  Polling int             `xml:"polling,attr"`
  Requests int            `xml:"requests,attr"`
  Hold int                `xml:"hold,attr"`
  From string             `xml:"from,attr"`
  Accept string           `xml:"accept,attr"`
  SID string              `xml:"sid,attr"`
  Secure bool             `xml:"secure,attr"`
  Charsets string         `xml:"charsets,attr"`
  RestartLogic bool       `xml:"xmpp:restartlogic,attr"`
  XMPPVersion string      `xml:"xmpp:version,attr"`
  AuthID string           `xml:"authid,attr"`
  Namespace string        `xml:"xmlns,attr"`
  XMPPNamespace string    `xml:"xmlns:xmpp,attr"`
  StreamNamespace string  `xml:"xmlns:stream,attr"`
  Features Feature        `xml:"stream:features"`
}

func NewResponse(r *Request) *Body {
  return &Body{
    Wait: r.Wait,
    Inactivity: 30,
    Polling: 5,
    Requests: 2,
    Hold: r.Hold,
    From: r.To,
    Accept: "deflate,gzip",
    SID: "thisismysessionid",
    Secure: false,
    Charsets: "ISO_8859-1 ISO-2022-JP",
    RestartLogic: true,
    XMPPVersion: r.XMPPVersion,
    AuthID: "ServerStreamID",
    Namespace: "http://jabber.org/protocol/httpbind",
    XMPPNamespace: "urn:xmpp:xbosh",
    StreamNamespace: "http://etherx.jabber.org/streams",
    Features: Feature{
      Bind: Bind{
        Namespace: "urn:ietf:params:xml:ns:xmpp-bind",
      },
    },
//     Features: Feature{
//       Mechanisms: Mechanisms {
//         Mechanisms: []string {
//           "PLAIN",
//         },
//       },
//     },
  }
}

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/http-bind", httpBindHandler)
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

func httpBindHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("\n")
  log.Printf("HTTP bind request incoming from %s\n\n", r.URL.Host + r.URL.Path)
  body, err := ioutil.ReadAll(r.Body)
  if err != nil || len(body) == 0 {
    log.Print(err)
    time.Sleep(3 * time.Second)
    log.Printf("Disconnecting\n")
    return
  }
  for k, v := range r.Header {
    fmt.Printf("%s: \"%s\"\n", k, v)
  }
  fmt.Printf("\n%s\n", body)
  req := Request{}
  err = xml.Unmarshal(body, &req)
  if err != nil {
    log.Print(err)
  }
  c := make(chan *Body, 1)
  go respondTo(&req, c)
  select {
  case resp := <-c:
    log.Printf("Success!!\n")
    bytes, err := xml.Marshal(struct {
     Body
     XMLName struct{}     `xml:"body"`
    }{Body: *resp})
    if err != nil {
      log.Print(err)
      return
    }
    _, err = w.Write(bytes)
    fmt.Printf("Response body:\n\n%s\n\n", string(bytes))
    if err != nil {
      log.Print(err)
    }

  case <-time.After(time.Duration(req.Wait) * time.Second):
    log.Printf("Request Timeout\n")
  }
}

func respondTo(r *Request, c chan *Body) {
  log.Printf("Message arrived for %s\n", r.To)
  c <-NewResponse(r)
}
