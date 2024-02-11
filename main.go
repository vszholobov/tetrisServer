package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tetrisServer/server"

	"github.com/gorilla/mux"
)

// https://github.com/gorilla/websocket/blob/main/examples/echo/server.go
func main() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/session", server.GetSessionsList)
	router.HandleFunc("/session/create", server.CreateSession)
	router.HandleFunc("/session/connect/{sessionId}", server.ConnectToSession)
	//http.HandleFunc("/", server.Home)
	log.Fatal(http.ListenAndServe(*server.Addr, router))
}
