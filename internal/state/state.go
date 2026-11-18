package state

import "github.com/AyuorusAguilar/Biscord/internal/db"

type State struct {
	ClientDir string
	Queries   *db.Queries
	K string
}
