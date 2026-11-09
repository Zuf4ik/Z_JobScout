package user

import "time"

type Session struct {
	ID        int64
	UserID    int64
	CreatedAt time.Time
	ExpiresAt time.Time
}
