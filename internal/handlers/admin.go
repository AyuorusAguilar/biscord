package handlers

import (
	"net/http"

	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func ResetAdmin(w http.ResponseWriter, r *http.Request, s state.State) {
	s.Queries.ResetDB(r.Context())
	logging.InnerLog("DATABASE RESET")
}