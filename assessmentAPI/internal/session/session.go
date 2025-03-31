package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Creates a new session manager for MongoDB.
func NewSessionManager(mongoURI string) (*SessionManager, error) {
	clientOptions := options.Client().ApplyURI(mongoURI).SetConnectTimeout(3 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database(DatabaseName).Collection(CollectionName)

	// TTL index on ExpiresAt field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}

	// index SessionID
	sessionIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "session_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(ctx, sessionIDIndex)
	if err != nil {
		return nil, err
	}

	return &SessionManager{
		client:     client,
		collection: collection,
	}, nil
}

// Creates a new user session and stores it in MongoDB.
func (sm *SessionManager) CreateSession(ctx context.Context) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	now := time.Now()
	session := Session{
		SessionID:   sessionID,
		Preferences: Preferences{PollingFrequency: 5, SortOrder: "name_alphabetical_ascending", HiddenDevices: []string{}, Visits: 0, ShowVisibilityControls: false},
		CreatedAt:   now,
		ExpiresAt:   now.Add(SessionExpirationTime),
	}

	opts := options.InsertOne()
	_, err = sm.collection.InsertOne(ctx, session, opts)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Retrieves the user session given the session ID from the client HTTPOnly cookie.
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	filter := bson.M{"session_id": sessionID}

	var session Session
	err := sm.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		sm.DeleteSession(ctx, sessionID)
		return nil, errors.New("session expired")
	}

	return &session, nil
}

// Updates the user preferences given the session ID and the complete user preferences object.
// Returns the an unordered bson document of the updated preferences.
func (sm *SessionManager) UpdateSession(ctx context.Context, sessionID string, preferences Preferences) bson.M {
	filter := bson.M{"session_id": sessionID}
	update := bson.M{
		"$set": bson.M{
			"preferences": preferences,
			"expires_at":  time.Now().Add(SessionExpirationTime),
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDocument bson.M
	err := sm.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			log.Fatal(err)
		}
	}

	return updatedDocument
}

// Deletes a user session given the session ID from the client HTTPOnly cookie.
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	filter := bson.M{"session_id": sessionID}
	_, err := sm.collection.DeleteOne(ctx, filter)
	return err
}

// Generates a unique session ID to be used as the value of the session cookie,
// as well as the key to retreiving the user's preferences from MongoDB.
func generateSessionID() (string, error) {
	b := make([]byte, SessionIDLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Retrieves the user session ID from the context.
func GetSessionFromContext(ctx context.Context) string {
	if sessionID, ok := ctx.Value(SessionCookieName).(string); ok {
		return sessionID
	}
	return ""
}

// Retrieves the user session data from the context.
func GetSessionDataFromContext(ctx context.Context) Preferences {
	if data, ok := ctx.Value("preferences").(Preferences); ok {
		return data
	}
	return Preferences{}
}
