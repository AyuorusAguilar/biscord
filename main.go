package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func main() {

	s := state.State{
		ClientDir: "client/",
	}

	os.Remove("log.txt")
	log("Booting...")

	mux := http.NewServeMux()
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(s.ClientDir+"css/"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(s.ClientDir+"js/"))))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(s.ClientDir+"img/"))))
	route(mux, s)
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log("Listening at Port 8080...")
	log("Check out http://localhost:8080/...")
	server.ListenAndServe()
}

func log(s string) {
	s = fmt.Sprintf("%s\n", s)
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s\n", "Logging failed!")
	}
	logFile.WriteString(s)
	fmt.Print(s)
}
