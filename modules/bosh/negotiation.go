package bosh

import (
  "github.com/yggie/EduChatSpike/modules/sasl"
  "github.com/yggie/EduChatSpike/modules/xmpp"
)

func createNewSession(r *Request) string {
  return `` +
  `<body wait="60" ` +
        `polling="5" ` +
        `inactivity="30" ` +
        `requests="2" ` +
        `hold="1" ` +
        `maxpause="120" ` +
        `sid="thisismysessionid" ` +
        `charsets="ISO_8859-1 ISO-2022-JP" ` +
        `ver="1.6" ` +
        `from="educhat.org" ` +
        `xmlns="http://jabber.org/protocol/httpbind">` +
      `<stream:features>` +
        `<mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">` +
          `<mechanism>SCRAM-SHA-1</mechanism>` +
        `</mechanisms>` +
      `</stream:features>` +
   `</body>`
}

func initSASL(r *Request) (string, error) {
  b, err := sasl.Init(r.Auth[0].Mechanism, r.SID, r.Auth[0].Data)
  if err != nil {
    return "", err
  }
  return `<body xmlns="http://jabber.org/protocol/httpbind">` + xmpp.Challenge(string(b)) + `</body>`, nil
}

func respondSASL(r *Request) (string, error) {
  b, err := sasl.Respond(r.SID, r.Response[0].Data)
  if err != nil {
    return "", err
  }
  return `<body xmlns="http://jabber.org/protocol/httpbind">` + xmpp.Success(string(b)) + `</body>`, nil
}

func restartSession(r *Request) (string, error) {
  if sasl.WasVerified(r.SID) {
    return `<body xmlns="http://jabber.org/protocol/httpbind" xmlns:stream="http://etherx.jabber.org/streams">` + xmpp.RestartStream() + `</body>`, nil
  }

  return "", InvalidSessionError
}
