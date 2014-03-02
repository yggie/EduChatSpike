package xmpp

type OnProcessListener interface {
  InitSASL(auth *Auth) error
  RespondSASL(clientResponse *Response) error
  Unprocessed() error
  Complete(response string) error
}

func (s *Stanza) Process(listener OnProcessListener) error {
  if s.HasAuth() {
    return listener.InitSASL(&s.Auth[0])

  } else if s.HasResponse() {
    return listener.RespondSASL(&s.Response[0])

  }

  response := ""
  if s.HasInfoQuery() {
    s, err := s.ProcessInfoQueries()
    if err != nil {
      return err
    }

    response += s
  }

  if response == "" {
    return listener.Unprocessed()
  } else {
    return listener.Complete(response)
  }
}

func (s *Stanza) ProcessInfoQueries() (string, error) {
  result := ""
  for _, iq := range s.Iq {
    str, err := s.ProcessInfoQuery(&iq)
    if err != nil {
      return "", err
    }
    result += str
  }

  return result, nil
}

func (s *Stanza) ProcessInfoQuery(iq *Iq) (string, error) {
  switch iq.Type {
  case "set":
    return `<iq id="` + iq.ID + `" type="result" xmlns="jabber:client"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><jid>testuser@examples.org/acompletelyrandomresource</jid></bind></iq>`, nil

  default:
    // TODO return error message
    return "", nil
  }
}
