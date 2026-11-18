package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/AyuorusAguilar/Biscord/internal/db"
	"github.com/AyuorusAguilar/Biscord/internal/handlers"
	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()
	// Database connection
	database, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		logging.Log("FATAL: Failed to open db connection!")
		os.Exit(1)
	}
	// State declaration
	s := state.State{
		ClientDir: "./client/",
		Queries: db.New(database),
		K: os.Getenv("K"),
	}
	
	// Log init
	os.Remove("log.txt")
	logging.Log("Booting...")
	// Multiplex initialization
	mux := http.NewServeMux()
	// FileServer routes
	mux.HandleFunc("/css/", logging.LogFileServerMiddleware("css", s))
	mux.HandleFunc("/js/",logging.LogFileServerMiddleware("js", s))
	mux.HandleFunc("/img/", logging.LogFileServerMiddleware("img", s))
	// Client routes
	route(mux, s)
	// Api routes
	mux.HandleFunc("/register", logging.LogApiMiddleware("/log", handlers.StateHandlerFuncMiddleware(s, handlers.RegisterUser)))
	// Server configuraiton
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// Server Start
	logging.Log("Listening at Port 8080...")
	logging.Log("Check out http://localhost:8080/...")
	server.ListenAndServe()
}


