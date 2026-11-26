package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AyuorusAguilar/Biscord/internal/db"
	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
	"github.com/AyuorusAguilar/Biscord/internal/token"
	"github.com/alexedwards/argon2id"
)

// RegisterUser is a handler that registers a user in the database.
// That's all it does though, to get refresh tokens and JWTs use
// the LoginUser handler
//
// It expects a request body in the following format:
// {
//	"username": "username",
//	"password": "userpass"
// }
func RegisterUser(w http.ResponseWriter, r *http.Request, s state.State) {
	// Read request
	type expectedRequestBody struct {
		Username	string `json:"username"`
		Pass		string `json:"password"`
	}
	var req expectedRequestBody
	if UnmarshallJson(r.Body, &req) != nil {
		logging.InnerLog("Error!")
		WriteError(w, 400, "Invalid Request Body")
		return
	}
	// Validate request
	if req.Pass == "" || req.Username == "" {
		logging.InnerLog("Error! One or more request fields are empty")
		WriteError(w, 400, "Invalid Request Body: One or more request fields are empty")
		return
	}
	
	exists, err := s.Queries.UsernameExists(r.Context(), req.Username)
	if exists {
		logging.InnerLog("Error! Username already in use")
		WriteError(w, 400, "Username is already taken")
		return
	}
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 500, "Unexpected Error: Error consulting username disponibility")
		return
	}
	// Data registration
	DMchannelID, err := s.Queries.CreateDMChannel(r.Context()) 
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 500, "Unexpected Error: Error creating DMchannel")
		return
	}
	
	passhash, err := argon2id.CreateHash(req.Pass, argon2id.DefaultParams)
	
	user, err := s.Queries.RegisterUser(r.Context(), db.RegisterUserParams{
		Username: req.Username,
		Pass: passhash,
		Dmchannelid: DMchannelID,
	})
	if err != nil {
		logging.InnerLog("Error registering user in database!")
		WriteError(w, 500, "Invalid Request Body: Error registering user in database!")
		return
	}

	
	// Success!
	logging.InnerLog(fmt.Sprintf("Succesfully registered user \"%s\"", user.Username))
	WriteJsonResponse(w, 201, map[string]string{
		"username": user.Username,
		"created_at": user.CreatedAt.String(),
	})
}

// LoginUser is a handler that creates an active session and returns
// a refresh token and a JWT, it also forces the last active session
// for that user to expire.
// 
// It expects a request body in the following format:
// {
//	"username": "username",
//	"password": "userpass"
// }
func LoginUser(w http.ResponseWriter, r *http.Request, s state.State) {
	// Read request
	type expectedRequestBody struct {
		Username	string `json:"username"`
		Pass		string `json:"password"`
	}
	var req expectedRequestBody
	if UnmarshallJson(r.Body, &req) != nil {
		logging.InnerLog("Error!")
		WriteError(w, 400, "Invalid Request Body")
		return
	}
	// Validate request
	if req.Pass == "" || req.Username == "" {
		logging.InnerLog("Error! One or more request fields are empty")
		WriteError(w, 400, "Invalid Request Body: One or more request fields are empty")
		return
	}
	
	userID, err := s.Queries.GetIdByUsername(r.Context(), req.Username)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 400, "Unexpected Error: Error consulting account existence")
		return
	}
	if last_session, err := s.Queries.GetLastSession(r.Context(), userID); err == nil {
		if last_session.ExpiresAt.After(time.Now()) {
			_, err := s.Queries.ExpireSession(r.Context(), last_session.ID)
			if err != nil {
				logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
				WriteError(w, 400, "Unexpected Error: Error expiring last unexpired session")
				return
			}
		}
	}
	// Creating Data
	refresh_token, refresh_token_hash, err := token.CreateRefreshToken()
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error creating refresh token: %v", err))
		WriteError(w, 500, "Unexpected Error: Error creating refresh token")
		return
	}
	jut, err := token.CreateJWT(int(userID), s.K, s.JWTDuration)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 500, "Unexpected Error: Error creating JWT")
		return
	}
	// Registering session
	session, err := s.Queries.RegisterSession(r.Context(), db.RegisterSessionParams{
		Userid: userID,
		RefreshTokenHash: refresh_token_hash[:],
		ExpiresAt: time.Now().Add(s.SessionDuration),
	})
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 500, "Unexpected Error: Error registrating session on database")
		return
	}
	
	// Success!
	logging.InnerLog(fmt.Sprintf("Succesfully login for userID \"%d\"", userID))
	WriteJsonResponse(w, 200, map[string]any{
		"refresh_token": base64.StdEncoding.EncodeToString(refresh_token),
		"jwt": jut,
		"logged_at": session.CreatedAt.String(),
		"expires_at": session.ExpiresAt.String(),
	})
}

// AuthRefreshUser is a handler that returns a JWT from a valid refresh
// token.
// 
// It expects a header in the following format:
// "Authorization: Bearer <refresh_token>"
func AuthRefreshUser(w http.ResponseWriter, r *http.Request, s state.State)  {
	// Token processing
	refreshToken, err := stripBearerFromToken(r.Header)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error! Invalid Authorization header: %s", r.Header.Get("Authorization")) )
		WriteError(w, 400, "Invalid Refresh Token")
		return
	}
	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error! Invalid Authorization header: %s", r.Header.Get("Authorization")) )
		WriteError(w, 400, "Invalid Refresh Token")
		return
	}
	refreshTokenHash := token.GetHash(refreshTokenBytes)
	// Token validation
	session, err := s.Queries.GetSessionByRefreshToken(r.Context(), refreshTokenHash[:])
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error consulting session on database: %v", err))
		WriteError(w, 400, "Session does not exist. Hit api/login to get a refresh token")
		return
	}
	if session.ManualExpiration {
		logging.InnerLog("Error! Session manually expired.")
		WriteError(w, 400, "Session was manually expired, hit api/login to get a new refresh toke")
		return
	}
	if session.ExpiresAt.Before(time.Now()) {
		logging.InnerLog("Error! Session expired.")
		WriteError(w, 400, "Session has expired, hit api/login to get a new refresh toke")
		return
	}
	// Success!
	jut, err := token.CreateJWT(int(session.Userid), s.K, s.JWTDuration)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 500, "Unexpected Error: Error creating JWT")
		return
	}
	logging.InnerLog(fmt.Sprintf("Succesfully JWT creation for userID \"%d\"", session.Userid))
	WriteJsonResponse(w, 200, map[string]any{
		"jwt": jut,
	})
}

// LogoutUser is a handler sets the manual_expiration
// field of a session entry (the one specified in the
// Authorization header) to false so no new JWT can
// be created from that refresh token.
// 
// It expects a header in the following format:
// "Authorization: Bearer <refresh-token>"
func LogoutUser(w http.ResponseWriter, r *http.Request, s state.State)  {
	// Token processing
	refreshToken, err := stripBearerFromToken(r.Header)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error! Invalid Authorization header: %s", r.Header.Get("Authorization")) )
		WriteError(w, 400, "Invalid Refresh Token")
		return
	}
	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error! Invalid Authorization header: %s", r.Header.Get("Authorization")) )
		WriteError(w, 400, "Invalid Refresh Token")
		return
	}
	refreshTokenHash := token.GetHash(refreshTokenBytes)
	// Token validation
	session, err := s.Queries.GetSessionByRefreshToken(r.Context(), refreshTokenHash[:])
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Error consulting session on database: %v", err))
		WriteError(w, 400, "Session does not exist.")
		return
	}
	if session.ManualExpiration {
		logging.InnerLog("Error! Session manually expired.")
		WriteError(w, 400, "Session was already manually expired.")
		return
	}
	if session.ExpiresAt.Before(time.Now()) {
		logging.InnerLog("Error! Session expired.")
		WriteError(w, 400, "Session has already expired")
		return
	}
	// Session closing!
	_, err = s.Queries.ExpireSession(r.Context(), session.ID)
	if err != nil {
		logging.InnerLog(fmt.Sprintf("Unexpected Error: %v", err))
		WriteError(w, 500, "Error closing session")
		return
	}
	// Success!
	logging.InnerLog(fmt.Sprintf("Succesfully JWT creation for userID \"%d\"", session.Userid))
	w.WriteHeader(200)
}

func stripBearerFromToken(h http.Header) (string, error) {
	auth := h.Get("Authorization")
	token, found := strings.CutPrefix(auth, "Bearer ")
	if !found {
		return "", fmt.Errorf("Invalid Format")
	}
	return token, nil
}

func getId(w http.ResponseWriter, r *http.Request, s state.State) (int64, error)  {
	jwt, err := stripBearerFromToken(r.Header)
	if err != nil {
		logging.InnerLog("Error gettin jwt from request")
		WriteError(w, 400, "Invalid JWT")
		return 0, fmt.Errorf("")
	}
	BearerID, err := token.ValidateJWT(jwt, s.K)
	if err != nil {
		logging.InnerLog("Error processing jwt")
		WriteError(w, 400, "Invalid JWT")
		return 0, fmt.Errorf("")
	}
	return BearerID, nil
}