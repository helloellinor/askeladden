package services

import (
	"sync"
)

// GlobalEmbedSessions holds active embed creation sessions
var GlobalEmbedSessions = &EmbedSessionManager{
	sessions: make(map[string]*EmbedSession),
}

// EmbedSessionManager manages embed creation sessions globally
type EmbedSessionManager struct {
	sessions map[string]*EmbedSession
	mutex    sync.RWMutex
}

// StartSession starts a new embed creation session for a user
func (esm *EmbedSessionManager) StartSession(userID, guildID string) {
	esm.mutex.Lock()
	defer esm.mutex.Unlock()

	esm.sessions[userID] = &EmbedSession{
		UserID:  userID,
		GuildID: guildID,
		Step:    0,
		Color:   ColorInfo, // Default color
	}
}

// GetSession gets an active session for a user
func (esm *EmbedSessionManager) GetSession(userID string) (*EmbedSession, bool) {
	esm.mutex.RLock()
	defer esm.mutex.RUnlock()

	session, exists := esm.sessions[userID]
	return session, exists
}

// RemoveSession removes a session for a user
func (esm *EmbedSessionManager) RemoveSession(userID string) {
	esm.mutex.Lock()
	defer esm.mutex.Unlock()

	delete(esm.sessions, userID)
}

// HasActiveSession checks if a user has an active embed creation session
func (esm *EmbedSessionManager) HasActiveSession(userID string) bool {
	esm.mutex.RLock()
	defer esm.mutex.RUnlock()

	_, exists := esm.sessions[userID]
	return exists
}
