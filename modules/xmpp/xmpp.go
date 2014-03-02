package xmpp

type Auth struct {
  Namespace string      `xml:"xmlns,attr"`
  Mechanism string      `xml:"mechanism,attr"`
  Data []byte           `xml:",chardata"`
  XMLName interface{}   `xml:"auth"`
}

func (a *Auth) Exists() bool {
  return len(a.Namespace) != 0
}

type Response struct {
  Namespace string      `xml:"xmlns,attr"`
  Data []byte           `xml:",chardata"`
  XMLName interface{}   `xml:"response"`
}

func (r *Response) Exists() bool {
  return len(r.Namespace) != 0
}

func Challenge(content string) string {
  return `<challenge xmlns="urn:ietf:params:xml:ns:xmpp-sasl">` + content + `</challenge>`
}

func Success(content string) string {
  return `<success xmlns="urn:ietf:params:xml:ns:xmpp-sasl">` + content + `</success>`
}

func RestartStream() string {
  return `<stream:features><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"/></stream:features>`
}
