package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
)

// GetUserById is a handler that returns a user object.
//
// It expects a "userid" declared in the pathvalue
func GetUserById(w http.ResponseWriter, r *http.Request, s state.State) {
	// Read request
	requestedUserId := r.PathValue("userid")
	userID, err := strconv.ParseInt(requestedUserId, 10, 64)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error parsing requested userID!: %s", requestedUserId))
		WriteError(w, 400, "Invalid userID")
		return
	}
	// Validate request
	user, err := s.Queries.GetUserById(r.Context(), userID)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error getting user from database: %v", err))
		WriteError(w, 400, "Invalid userID")
		return
	}
	// Success!
	logging.InnerLog(fmt.Sprintf("Succesfully delivered user \"%s\"", user.Username))
	WriteJsonResponse(w, 201, map[string]any{
		"username": user.Username,
		"profilepic": user.Profilepic,
		"DMchannelID": user.Dmchannelid,
	})
}