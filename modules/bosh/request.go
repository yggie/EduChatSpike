package bosh

import (
  "fmt"
  "time"
  "net/http"
  "io/ioutil"
  "encoding/xml"
)

// The session creation request body as defined by XEP-0206
type Request struct {

  // embeds the standard XMPP message body
  InnerXML string       `xml:",innerxml"`

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
  Wait  time.Duration   `xml:"wait,attr"`

  // The default language of the XML character data
  Lang  string          `xml:"xml:lang,attr"`

  XMPPVersion string    `xml:"xmpp:version,attr"`

  Namespace string      `xml:"xmlns,attr"`

  XMPPNamespace  string `xml:"xmlns:xmpp,attr"`

  // highest version of the BOSH protocol that the client supports, as defined in XEP-0124
  Version string        `xml:"ver,attr"`

  // The session ID, present only for existing sessions
  SessionID string      `xml:"sid,attr"`

  XMLName xml.Name      `xml:"body"`
}

func (r *Request) ShouldRestart() bool {
  return len(r.Restart) != 0
}

func (r *Request) HasSessionID() bool {
  return r.SessionID != ""
}

func ParseRequest(r *http.Request) (*Request, error) {
  // read body data from buffer
  body, err := ioutil.ReadAll(r.Body)
  if err != nil || len(body) == 0 {
    return nil, err
  }

  // debug lines
  fmt.Printf("\nRequest:\n%s\n", body)

  // parse request body
  req := Request{}
  err = xml.Unmarshal(body, &req)
  if err != nil {
    return nil, err
  }

  return &req, nil
}

