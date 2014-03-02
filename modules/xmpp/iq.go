package xmpp

import (
  "encoding/xml"
)

type Iq struct {
  ID string             `xml:"id,attr"`
  Type string           `xml:"type,attr"`
  Bind []Bind           `xml:"bind"`
  Ping []Ping           `xml:"ping"`
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

type Ping struct {
  Namespace string      `xml:"xmlns,attr"`
  XMLName xml.Name      `xml:"ping"`
}

type Resource struct {
  Text string           `xml:",chardata"`
}

func (iq *Iq) HasPing() bool {
  return len(iq.Ping) != 0
}
