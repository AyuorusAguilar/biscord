package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AyuorusAguilar/Biscord/internal/logging"
	"github.com/AyuorusAguilar/Biscord/internal/state"
)

func route(mux *http.ServeMux, s state.State) {
	routes := openDict()
	for route, filepath := range routes {
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, fmt.Sprintf("%s%s", s.ClientDir, filepath))
			logging.Log(fmt.Sprintf("Served %s file", filepath))
		})
	}
}

func openDict() map[string]string {
	routeFile, err := os.OpenFile("routes.json", os.O_RDONLY, 0644)
	if err != nil {
		logging.Log("An error ocurred while opening the routeFile!. Exiting...")
		os.Exit(1)
	}

	routeBytes, err := io.ReadAll(routeFile)
	if err != nil {
		logging.Log("An error ocurred while reading the routeFile!. Exiting...")
		os.Exit(1)
	}

	var routeMap map[string]string
	
	err = json.Unmarshal(routeBytes, &routeMap)
	if err != nil {
		logging.Log("An error ocurred while unmarshalling the routeFile!. Exiting...")
		os.Exit(1)
	}
	
	return routeMap
}
