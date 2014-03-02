package xmpp

import (
  "encoding/xml"
)

type Iq struct {
  ID string             `xml:"id,attr"`
  Type string           `xml:"type,attr"`
  Bind Bind
  XMLName xml.Name      `xml:"iq"`
}

type Bind struct {
  Namespace string      `xml:"xmlns,attr"`
  Resource Resource
  XMLName xml.Name      `xml:"bind"`
}

type JID struct {
  Text string           `xml:",chardata"`
  XMLName xml.Name      `xml:"jid"`
}

type Resource struct {
  Text string           `xml:",chardata"`
}
