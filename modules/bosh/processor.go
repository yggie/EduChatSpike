package bosh

import (
  "fmt"
  "time"
  "net/http"
  "github.com/yggie/EduChatSpike/modules/xmpp"
  "github.com/yggie/EduChatSpike/modules/sasl"
)

type RequestProcessor struct {
  Request *Request
  Writer http.ResponseWriter
}

// initializes an SASL handshake
func (r *RequestProcessor) InitSASL(auth *xmpp.Auth) error {
  response, data, err := sasl.Init(auth.Mechanism, auth.Data)
  if err != nil {
    return err
  }

  s, err := GetSession(r.Request.SessionID)
  if err != nil {
    return err
  }

  s.Clear()
  s.Mechanism = auth.Mechanism
  s.Data = data

  r.WriteResponse(`<body xmlns="http://jabber.org/protocol/httpbind">` + xmpp.Challenge(string(response)) + `</body>`)

  return nil
}

// response to an SASL handshake
func (r *RequestProcessor) RespondSASL(clientResponse *xmpp.Response) error {
  session, err := GetSession(r.Request.SessionID)
  if err != nil {
    return err
  }

  response, err := sasl.Respond(session.Mechanism, clientResponse.Data, session.Data)
  if err != nil {
    return err
  }
  session.Clear()
  session.Verified = true

  r.WriteResponse(`<body xmlns="http://jabber.org/protocol/httpbind">` + xmpp.Success(string(response)) + `</body>`)

  return nil
}

func (r *RequestProcessor) Unprocessed() error {
  var err error
  if r.Request.SessionID == "" {
    err = r.CreateNewSession()

  } else if r.Request.ShouldRestart() {
    err = r.RestartSession()

  } else {
    time.Sleep(60 * time.Second)
  }

  return err
}

func (r *RequestProcessor) Complete(response string) error {
  r.WriteResponse(`<body xmlns="http://jabber.org/protocol/httpbind">` + response + `</body>`)

  return nil
}

// initializes a new session
func (r *RequestProcessor) CreateNewSession() error {
  session := NewSession()
  // TODO actually read request data to determine settings
  r.WriteResponse(`` +
  `<body wait="60" ` +
        `polling="5" ` +
        `inactivity="30" ` +
        `requests="2" ` +
        `hold="1" ` +
        `maxpause="120" ` +
        `sid="` + session.ID + `" ` +
        `charsets="ISO_8859-1 ISO-2022-JP" ` +
        `ver="1.6" ` +
        `from="educhat.org" ` +
        `xmlns="http://jabber.org/protocol/httpbind">` +
      `<stream:features>` +
        `<mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">` +
          `<mechanism>SCRAM-SHA-1</mechanism>` +
        `</mechanisms>` +
      `</stream:features>` +
   `</body>`)

   // could be used in the future
   return nil
}

func (r *RequestProcessor) RestartSession() error {
  session, err := GetSession(r.Request.SessionID)
  if err != nil {
    return err
  }

  if session.Verified {
    r.WriteResponse(`<body xmlns="http://jabber.org/protocol/httpbind" xmlns:stream="http://etherx.jabber.org/streams">` + xmpp.RestartStream() + `</body>`)
  } else {
    return InvalidSessionError
  }

  return nil
}

func (r *RequestProcessor) WriteResponse(response string) {
  // set up default reply headers
  r.Writer.Header().Set("Content-Type", "text/xml")

  fmt.Printf("\nResponse:\n%s\n\n", response)
  r.Writer.Write([]byte(response))
}

