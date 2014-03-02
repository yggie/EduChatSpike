package bosh

import (
  "errors"
)

var (
  sessionNotFound = errors.New("Session not found")

  sessionStore = make(map[string]*Session)
)

type Session struct {
  ID string
  Data []byte
  Mechanism string
  Verified bool
}

func NewSession() *Session {
  for {
    sid := "thisismysessionid"
    _, found := sessionStore[sid]

    if !found {
      session := &Session{ ID: sid }
      sessionStore[sid] = session
      return session
    }

    // TODO generate unique session ids
    return sessionStore["thisismysessionid"]
  }
}

func GetSession(sessionID string) (*Session, error) {
  session, found := sessionStore[sessionID]
  if !found {
    return nil, sessionNotFound
  }

  return session, nil
}

func (s *Session) Clear() {
  s.Data = nil
}
