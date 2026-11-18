package logging

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func Log(s string) {
	s = fmt.Sprintf("%s\n", s)
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s\n", "Logging failed!")
	}
	logFile.WriteString(s)
	fmt.Print(s)
}

func LogFileServerMiddleware(dir string, s state.State) http.HandlerFunc{
	actuaFunc := http.StripPrefix(fmt.Sprintf("/%s/", dir), http.FileServer(http.Dir(s.ClientDir+dir+"/")))
	return func(w http.ResponseWriter, r *http.Request) {
		Log(fmt.Sprintf("Served a file from /%s/", dir))
		actuaFunc.ServeHTTP(w, r)
	}
}

func LogApiMiddleware(endpoint string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Log(fmt.Sprintf("Recived call to %s endpoint", endpoint))
		handlerFunc.ServeHTTP(w, r)
	}
}
// 