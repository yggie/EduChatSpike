package xmpp

const (
  DOMAIN = "educhat.spike"
)

var (
  activeUsers = make(map[string]*ActiveUser)
)

type Response struct {
  Namespace string      `xml:"xmlns,attr"`
  Data []byte           `xml:",chardata"`
  XMLName interface{}   `xml:"response"`
}

func (r *Response) Exists() bool {
  return len(r.Namespace) != 0
}

func RestartStream() string {
  return `<stream:features><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"/></stream:features>`
}

func NewResourceForUser(name string, receiver Receiver) *UserResource {
  user, ok := activeUsers[name]
  if !ok {
    // TODO verify user exists?
    user = &ActiveUser{ Name: name }
  }

  resource := UserResource{
    ID: "randomly-generated-id",
    User: user,
    Receiver: receiver,
  }

  user.Resources = append(user.Resources, resource)

  return &resource
}
