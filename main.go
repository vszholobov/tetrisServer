package main

import (
	"flag"
	"log"
	"net/http"
	"tetrisServer/server"
)

// https://github.com/gorilla/websocket/blob/main/examples/echo/server.go
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", server.Echo)
	http.HandleFunc("/", server.Home)
	log.Fatal(http.ListenAndServe(*server.Addr, nil))
}
