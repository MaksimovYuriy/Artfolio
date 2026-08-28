package entity

import "time"

type Session struct {
	ID        int64
	ActorID   int64
	Token     string
	ExpiresAt time.Time
}

type AuthenticatedSession struct {
	ID      int64
	ActorID int64
}
