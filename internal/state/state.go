package state

import (
	"time"

	"github.com/AyuorusAguilar/Biscord/internal/db"
)

type State struct {
	ClientDir string
	Queries   *db.Queries
	K string
	JWTDuration time.Duration
	SessionDuration time.Duration
}
