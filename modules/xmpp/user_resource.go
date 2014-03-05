package xmpp

type ActiveUser struct {
  Name string
  Resources []UserResource
}

type UserResource struct {
  ID string
  User *ActiveUser
  Receiver Receiver
}

type Receiver interface {
  Receive(message string)
}

func (r *UserResource) Send(message string) error {
  stanza, err := ParseMessage(message)
  if err != nil {
    return err
  }

  response := ""
  if stanza.HasInfoQuery() {
    for _, iq := range stanza.Iq {
      response += r.ProcessInfoQuery(iq)
    }
  }

  if len(response) != 0 {
    r.Push(response)
  }

  return nil
}

func (r *UserResource) Push(message string) {
  r.Receiver.Receive(message)
}

func (r *UserResource) UserName() string {
  return r.User.Name
}

func (r *UserResource) JID() string {
  return r.User.Name + "@" + DOMAIN + "/" + r.ID
}

func (r *UserResource) ProcessInfoQuery(iq Iq) string {
  response := ""
  if iq.HasPing() {
    if iq.Ping[0].Namespace == "urn:xmpp:ping" && iq.To == DOMAIN {
      response += "<iq id='" + iq.ID + "' type='result' to='" + r.JID() + "' from='" + DOMAIN + "' />"
    }
  }

  return response
}
