package bosh

import (
  "fmt"
  "github.com/yggie/EduChatSpike/modules/xmpp"
)

func parseData(r *Request) (string, error) {
  response := `<body xmlns="http://jabber.org/protocol/httpbind">`

  for _, iq := range r.Iq {
    iqres, err := respondToIq(&iq)
    if err != nil {
      return "", err
    }
    response += iqres
  }

  return response + "</body>", nil
}

func respondToIq(iq *xmpp.Iq) (string, error) {
  switch iq.Type {
  case "set":
    return `<iq id="` + iq.ID + `" type="result" xmlns="jabber:client"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><jid>testuser@examples.org/acompletelyrandomresource</jid></bind></iq>`, nil

  default:
    return "", fmt.Errorf("unsupported info query\n")
  }
}
