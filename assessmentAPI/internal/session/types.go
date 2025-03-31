package session

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Preferences struct {
	SortOrder              string   `bson:"sort_order" json:"sort_order"`
	HiddenDevices          []string `bson:"hidden_devices" json:"hidden_devices"`
	Visits                 int      `bson:"visits" json:"visits"`
	ShowVisibilityControls bool     `bson:"show_visibility_controls" json:"show_visibility_controls"`
	PollingFrequency       float64  `bson:"polling_frequency" json:"polling_frequency"`
}

type Session struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	SessionID   string        `bson:"session_id"`
	Preferences Preferences   `bson:"preferences"`
	CreatedAt   time.Time     `bson:"created_at"`
	ExpiresAt   time.Time     `bson:"expires_at"`
}

type SessionManager struct {
	client     *mongo.Client
	collection *mongo.Collection
}
