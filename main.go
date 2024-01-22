package main

import (
	"flag"
	"log"
	"net/http"
	"tetrisServer/server"

	"github.com/gorilla/mux"
)

// https://github.com/gorilla/websocket/blob/main/examples/echo/server.go
func main() {
	flag.Parse()
	log.SetFlags(0)
	router := mux.NewRouter()
	router.HandleFunc("/session", server.GetSessionsList)
	router.HandleFunc("/session/create", server.CreateSession)
	router.HandleFunc("/session/connect/{sessionId}", server.ConnectToSession)
	//http.HandleFunc("/", server.Home)
	log.Fatal(http.ListenAndServe(*server.Addr, router))
}
