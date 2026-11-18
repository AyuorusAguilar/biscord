package handlers

import (
	"net/http"

	"github.com/AyuorusAguilar/Biscord/internal/db"
	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func RegisterUser(w http.ResponseWriter, r *http.Request, s state.State) {
	user, err := s.Queries.RegisterUser(r.Context(), db.RegisterUserParams{
		Email: "ross@gmail.com",
		Pass: "ContraseniaSegura69",
		Username: "Ross",
	})
	if err != nil {
		logging.Log("%tError registering user!")
		WriteError(w, http.StatusNotModified, "Invalid Request Body")
		return
	}
	WriteJsonResponse(w, 201, map[string]string{
		"Email": user.Email,
		"Username": user.Username,
		"Created_at": user.CreatedAt.String(),
	})
}