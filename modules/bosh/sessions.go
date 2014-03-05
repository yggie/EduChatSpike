package bosh

import (
  "time"
  "errors"
  "github.com/yggie/EduChatSpike/modules/xmpp"
)

var (
  sessionNotFound = errors.New("Session not found")

  activeSessions = make(map[string]*Session)
)

type Session struct {
  ID string
  Wait time.Duration
  LastRequestID int64
  Connection xmpp.Connection
  MessageQueue chan string
}

func (s *Session) Bind(conn xmpp.Connection) {
  s.Connection = conn
}

func (s *Session) Receive(message string) {
  s.MessageQueue <-message
}

func CreateSession(sessionID string) *Session {
  session := Session{
    ID: sessionID,
    MessageQueue: make(chan string, 1),
    Wait: ServerSettings.Wait,
  }

  session.Connection = &xmpp.Anonymous{
    Receiver: &session,
    Binder: &session,
  }

  return &session
}

func NewSession() *Session {
  for {
    sid := "thisismysessionid"
    _, found := activeSessions[sid]

    if !found {
      session := CreateSession(sid)
      activeSessions[sid] = session
      return session
    }

    // TODO generate unique session ids
    return activeSessions["thisismysessionid"]
  }
}

func GetSession(sessionID string) (*Session, error) {
  session, found := activeSessions[sessionID]
  if !found {
    return nil, sessionNotFound
  }

  return session, nil
}

