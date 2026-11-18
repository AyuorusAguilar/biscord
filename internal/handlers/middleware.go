package handlers

import (
	"net/http"

	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func StateHandlerFuncMiddleware(s state.State, actualHandler func(w http.ResponseWriter, r *http.Request, s state.State)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actualHandler(w, r, s)
	}
}