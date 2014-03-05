package bosh

import (
  "fmt"
  "log"
  "time"
  "strconv"
  "net/http"
)

var (
  InvalidSessionError = fmt.Errorf("invalid session id")
  ExpiredRequest = fmt.Errorf("expired request")
  ServerSettings = Settings{
    Wait: 60,
    Name: "tardy-tanzanine",
  }
)

type Settings struct {
  Wait time.Duration
  Name string
}

func HttpBindHandler(w http.ResponseWriter, r *http.Request) {
  // parse request body
  request, err := ParseRequest(r)

  var response string
  if err == nil {

    response, err = HandleRequest(request)

    if err == ExpiredRequest {
      return
    }
  }

  // set default headers
 w.Header().Set("Content-Type", "text/xml")

  // debug line
  fmt.Printf("\nResponse:\n%s\n\n", response)

  if err != nil {
    log.Println(err)
    w.WriteHeader(400)
  }
  w.Write([]byte(response))
}

func HandleRequest(r *Request) (string, error) {
  if r.HasSessionID() {
    session, err := GetSession(r.SessionID)
    if err != nil {
      return "", err
    }

    if r.RID <= session.LastRequestID {
      return "", ExpiredRequest
    }

    session.LastRequestID = r.RID

    if r.ShouldRestart() {
      return RestartStream(session), nil

    } else {
      err = session.Connection.Send(r.InnerXML)
      if err != nil {
        return "", err
      }

      select {
      case response := <-session.MessageQueue:
        return "<body xmlns='http://jabber.org/protocol/httpbind'>" + response + "</body>", nil

        // time out
      case <-time.After(session.Wait * time.Second):
        // TODO return a non-blank response?
        return "", nil
      }
    }

  } else {
    return InitializeNewSession(), nil
  }
}

func InitializeNewSession() string {
  session := NewSession()
  return "" +
  "<body wait='" + strconv.Itoa(int(session.Wait)) + "' " +
        "polling='5' " +
        "inactivity='30' " +
        "requests='2' " +
        "hold='1' " +
        "maxpause='120' " +
        "sid='" + session.ID + "' " +
        "charsets='ISO_8859-1 ISO-2022-JP' " +
        "ver='1.6' " +
        "from='educhat.spike' " +
        "xmlns='http://jabber.org/protocol/httpbind'>" +
    "<stream:features>" +
      "<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
        "<mechanism>SCRAM-SHA-1</mechanism>" +
      "</mechanisms>" +
    "</stream:features" +
  "</body>"
}

func RestartStream(session *Session) string {
  return "<body xmlns='http://jabber.org/protocol/httpbind' xmlns:stream='http://etherx.jabber.org/streams'><stream:features><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></stream:features></body>"
}
