package api

import (
	"ProxyPool/storage"
	"ProxyPool/util"
	"encoding/json"
	"log"
	"net/http"
)

// VERSION for this program
const VERSION = "/v1"

// Run for request
func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc(VERSION+"/ip", ProxyHandler)
	log.Println("Starting server", util.NewConfig().Host)
	http.ListenAndServe(util.NewConfig().Host, mux)
}

// ProxyHandler .
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(storage.ProxyRandom())
		if err != nil {
			return
		}
		w.Write(b)
	}
}
