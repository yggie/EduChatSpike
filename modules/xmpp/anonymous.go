package xmpp

import (
  "log"
  "github.com/yggie/EduChatSpike/modules/sasl"
)

type Anonymous struct {
  Username string
  AuthMechanism string
  SavedData []byte
  Receiver Receiver
  Binder Binder
  Verified bool
}

type Binder interface {
  Bind(conn Connection)
}

func (a *Anonymous) Clear() {
  a.SavedData = nil
  a.AuthMechanism = ""
}

func (a *Anonymous) Send(message string) error {
  stanza, err := ParseMessage(message)
  if err != nil {
    return err
  }

  if a.Verified {
    if stanza.HasInfoQuery() {
      if stanza.Iq[0].Type == "set" &&
          stanza.Iq[0].ShouldBind() &&
          a.Verified &&
          a.Username != "" {
        resource := NewResourceForUser(a.Username, a.Receiver)
        a.Binder.Bind(resource)
        resource.Push("<iq id='" + stanza.Iq[0].ID + "' type='result' xmlns='jabber:client'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><jid>" + resource.JID() + "</jid></bind></iq>")
        return nil
      }

    }

  } else {
    if stanza.HasAuth() {
      a.InitSASL(&stanza.Auth[0])
      return nil

    } else if stanza.HasResponse() {
      a.RespondSASL(&stanza.Response[0])
      return nil

    }
  }

  // TODO error for unhandled message
  return nil
}

func (a *Anonymous) Push(message string) {
  a.Receiver.Receive(message)
}

func (a *Anonymous) UserName() string {
  return a.Username
}

func (a *Anonymous) JID() string {
  return ""
}

func (a *Anonymous) InitSASL(auth *Auth) {
  response, data, err := sasl.Init(auth.Mechanism, auth.Data)
  if err != nil {
    // TODO push error response
    log.Println(err)
    return
  }

  a.AuthMechanism = auth.Mechanism
  a.SavedData = data
  a.Username, err = auth.GetUserName()

  if err != nil {
    log.Println(err)
    return
  }

  a.Push("<challenge xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" + string(response) + "</challenge>")
}

func (a *Anonymous) RespondSASL(clientResponse *Response) {
  response, err := sasl.Respond(a.AuthMechanism, clientResponse.Data, a.SavedData)
  if err != nil {
    // TODO push error response
    log.Println(err)
    return
  }

  a.Clear()
  a.Verified = true

  a.Push("<success xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" + string(response) + "</success>")
}

