package models

import "time"

// These models match what sqlc will generate
// You can also just use the sqlc-generated types directly

type UserView struct {
	ID       int64
	Username string
}

type HaikuView struct {
	ID           int64
	FromUserID   int64
	ToUserID     int64
	Line1        string
	Line2        string
	Line3        string
	CreatedAt    time.Time
	FromUsername string
}
