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
		JWTDuration: 5 * time.Minute,
		SessionDuration: 30 * time.Minute,
	}
	
	// Enable *ON DELETE CASCADE* behavior
	_, err = database.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		logging.InnerLog("Error executing PRAGMA querie")
		os.Exit(1)
	}
	// Log init
	os.Remove("log.txt")
	logging.Log("Booting...")
	// Multiplex initialization
	mux := http.NewServeMux()
	// FileServer routes
	mux.HandleFunc("/css/",	logging.LogFileServerMiddleware("css", s))
	mux.HandleFunc("/js/",	logging.LogFileServerMiddleware("js", s))
	mux.HandleFunc("/img/", logging.LogFileServerMiddleware("img", s))
	// Client routes
	route(mux, s)
	// Api routes
	//		Authentication
	mux.HandleFunc("/register",		logging.LogApiMiddleware("/register", handlers.StateHandlerFuncMiddleware(s, handlers.RegisterUser)))
	mux.HandleFunc("/login",		logging.LogApiMiddleware("/login", handlers.StateHandlerFuncMiddleware(s, handlers.LoginUser)))
	mux.HandleFunc("/auth",			logging.LogApiMiddleware("/auth", handlers.StateHandlerFuncMiddleware(s, handlers.AuthRefreshUser)))
	mux.HandleFunc("/logout",		logging.LogApiMiddleware("/login", handlers.StateHandlerFuncMiddleware(s, handlers.LogoutUser)))
	//		Other (?)
	mux.HandleFunc("GET /user/{userid}",	logging.LogApiMiddleware("GET /user/{userid}", handlers.StateHandlerFuncMiddleware(s, handlers.GetUserById)))
	mux.HandleFunc("POST /messages",	logging.LogApiMiddleware("POST /messages", handlers.StateHandlerFuncMiddleware(s, handlers.PutMessage)))
	mux.HandleFunc("GET /messages",	logging.LogApiMiddleware("GET /messages", handlers.StateHandlerFuncMiddleware(s, handlers.GetMessagesFromChannelByID)))

	// Admin
	mux.HandleFunc("/reset",		logging.LogApiMiddleware("/login", handlers.StateHandlerFuncMiddleware(s, handlers.ResetAdmin)))
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
	logging.Log("Check out http://localhost:8080/...\nLog:\n")
	server.ListenAndServe()
}


