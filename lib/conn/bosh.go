package conn

import (
  "fmt"
  "log"
  "time"
  "yggie/EduChatSpike/lib/auth"
  "net/http"
  "io/ioutil"
  "encoding/xml"
)

type BOSHAuthResponse struct {
  Namespace string      `xml:"xmlns,attr"`
  Data []byte           `xml:",chardata"`
}

func (b *BOSHAuthResponse) Exists() bool {
  return len(b.Namespace) != 0
}

type BOSHAuth struct {
  Namespace string      `xml:"xmlns,attr"`
  Mechanism string      `xml:"mechanism,attr"`
  Data []byte           `xml:",chardata"`
}

func (b *BOSHAuth) Exists() bool {
  return len(b.Namespace) != 0
}

/// The session creation request body as defined by XEP-0206
type BOSHRequest struct {
  // Specifies the HTTP Content-Type of following responses
  Content string        `xml:"content,attr"`
  // The JabberID of the connecting entity
  From string           `xml:"from,attr"`
  // Maximum number of requests that can be alive at one time during the session
  Hold  int             `xml:"hold,attr"`
  // The request ID
  RID int64             `xml:"rid,attr"`
  // Target domain of the first stream
  To  string            `xml:"to,attr"`
  // For proxy sessions with other domains
  Route string          `xml:"route,attr"`
  Restart string        `xml:"restart,attr"`
  Secure bool           `xml:"secure,attr"`
  // Wait time in seconds
  Wait  int             `xml:"wait,attr"`
  // The default language of the XML character data
  Lang  string          `xml:"xml:lang,attr"`
  XMPPVersion string    `xml:"xmpp:version,attr"`
  Namespace string      `xml:"xmlns,attr"`
  XMPPNamespace  string `xml:"xmlns:xmpp,attr"`
  // highest version of the BOSH protocol that the client supports, as defined in XEP-0124
  Version string        `xml:"ver,attr"`
  // The session ID, present only for existing sessions
  SID string            `xml:"sid,attr"`
  // client-preferred authentication mechanism
  Auth BOSHAuth         `xml:"auth"`
  AuthResponse BOSHAuthResponse `xml:"response"`
}

func (r *BOSHRequest) ShouldRestart() bool {
  return len(r.Restart) != 0
}

func HttpBindHandler(w http.ResponseWriter, r *http.Request) {
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
  req := BOSHRequest{}
  err = xml.Unmarshal(body, &req)
  if err != nil {
    log.Print(err)
  }

  if req.SID == "" {
    err := createNewSession(w, &req)
    if err != nil {
      log.Print(err)
      return
    }
  } else {
    bytes, err := respondTo(&req)
    if err != nil {
      log.Println(err)
      return
    }
    _, err = w.Write(bytes)
  }
}

func createNewSession(w http.ResponseWriter, r *BOSHRequest) error {
  log.Printf("Initiating new session for arrived for %s\n", r.To)
  resp := `<body wait="60"
                 polling="5"
                 inactivity="30"
                 requests="2"
                 hold="1"
                 maxpause="120"
                 sid="thisismysessionid"
                 charsets="ISO_8859-1 ISO-2022-JP"
                 ver="1.6"
                 from="educhat.org"
                 xmlns="http://jabber.org/protocol/httpbind">
              <stream:features>
                <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
                  <mechanism>SCRAM-SHA-1</mechanism>
                </mechanisms>
              </stream:features>
           </body>`
  fmt.Printf("\n\nNew Session Created:\n\n%s\n", resp)
  _, err := w.Write([]byte(resp))
  return err
}

func respondTo(r *BOSHRequest) ([]byte, error) {
  if r.Auth.Exists() {
    b, err := auth.Init(r.Auth.Mechanism, r.SID, r.Auth.Data)
    if b != nil {
      bytes := []byte(`
      <body xmlns="http://jabber.org/protocol/httpbind">
        <challenge xmlns="urn:ietf:params:xml:ns:xmpp-sasl">` +
        string(b) +
       `</challenge>
      </body>`)
      return bytes, err
    }
    return nil, err
  } else if r.AuthResponse.Exists() {
    time.Sleep(60 * time.Second)
    bytes, err := auth.Respond(r.SID, r.AuthResponse.Data)
    return bytes, err
  } else if r.ShouldRestart() {
    time.Sleep(60 * time.Second)
  } else {
    time.Sleep(60 * time.Second)
  }
  return nil, nil
}

