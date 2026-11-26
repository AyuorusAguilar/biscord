package handlers

import (
	"fmt"
	"net/http"

	"github.com/AyuorusAguilar/Biscord/internal/db"
	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func PutMessage(w http.ResponseWriter, r *http.Request, s state.State) {
	type expectedRequestBody struct {
		Channel	int64 `json:"channel"`
		Content	string `json:"content"`
	}
	var req expectedRequestBody
	if UnmarshallJson(r.Body, &req) != nil {
		logging.InnerLog("Error!")
		WriteError(w, 400, "Invalid Request Body")
		return
	}
	// Validate request
	if req.Content == "" {
		logging.InnerLog("Error! One or more request fields are empty")
		WriteError(w, 400, "Invalid Request Body: One or more request fields are empty")
		return
	}
	userID, err := getId(w, r, s)
	if err != nil {
		return /* getId already handled the error */
	}
	_, err = s.Queries.SendMessage(r.Context(), db.SendMessageParams{
		Senderid: userID,
		Channelid: req.Channel,
		Content: req.Content,
	})
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error registering message on database: %v", err))
		WriteError(w, 500, "Unexpected error")
		return
	}
	// Success
	w.WriteHeader(200)
}
func GetMessagesFromChannelByID(w http.ResponseWriter, r *http.Request, s state.State) {
	type expectedRequestBody struct {
		Channel	int64 `json:"channel"`
		RequestedUserID	int64 `json:"userid"`
	}
	var req expectedRequestBody
	UnmarshallJson(r.Body, &req)
	// Validate request
	if req.Channel == 0 || req.RequestedUserID == 0 {
		logging.InnerLog("Error! One or more request fields are empty")
		WriteError(w, 400, "Invalid Request Body: One or more request fields are empty")
		return
	}
	userID, err := getId(w, r, s)
	if err != nil {
		return /* getId already handled the error */
	}
	logging.InnerLog(fmt.Sprintf("Request comming from userid: %d", userID))
	msgs, err := s.Queries.GetMessagesFromChannelByUserID(r.Context(), db.GetMessagesFromChannelByUserIDParams{
		Senderid: req.RequestedUserID,
		Channelid: req.Channel,
	})
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error registering message on database: %v", err))
		WriteError(w, 500, "Unexpected error")
		return
	}
	// Success
	logging.InnerLog(fmt.Sprintf("messages:\n%v", msgs))
	WriteJsonResponse(w, 200, msgs)
}