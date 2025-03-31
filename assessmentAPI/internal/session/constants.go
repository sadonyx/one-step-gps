package session

import "time"

const (
	// The byte length of the session ID before base64 encoding.
	SessionIDLength int = 32
	// The default expiration time for sessions.
	SessionExpirationTime = 72 * time.Hour
	// The name of the cookie storing the session ID.
	SessionCookieName string = "anonymous_session"
	// The name of the MongoDB database.
	DatabaseName string = "osgps_user_sessions"
	// The name of the MongoDB collection.
	CollectionName string = "user_sessions"
)
