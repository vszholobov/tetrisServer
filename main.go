package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tetrisServer/server"
)

// https://github.com/gorilla/websocket/blob/main/examples/echo/server.go
func main() {
	flag.Parse()
	log.SetFlags(0)
	router := mux.NewRouter()
	router.HandleFunc("/create", server.CreateSession)
	router.HandleFunc("/connect/{sessionId}", server.ConnectToSession)
	//http.HandleFunc("/", server.Home)
	log.Fatal(http.ListenAndServe(*server.Addr, router))
}
