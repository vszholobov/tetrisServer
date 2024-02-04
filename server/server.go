package server

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var Addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var Sessions = make(map[int64]*GameSession)

type CreateSessionResponse struct {
	SessionId int64 `json:"sessionId"`
}

func GetSessionsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Sessions)
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	gameSession := MakeGameSession()
	Sessions[gameSession.GetSessionId()] = gameSession
	response := CreateSessionResponse{SessionId: gameSession.GetSessionId()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ConnectToSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId, _ := strconv.ParseInt(vars["sessionId"], 10, 64)
	session := Sessions[sessionId]

	if session.Started {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Session already started"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	if session.FirstPlayerSession == nil {
		firstPlayerPieceGenerator := rand.New(rand.NewSource(sessionId))
		firstPlayerSession := MakePlayerSession(conn, firstPlayerPieceGenerator, session)
		session.FirstPlayerSession = firstPlayerSession
		log.Printf("Session %d created", sessionId)
	} else {
		secondPlayerPieceGenerator := rand.New(rand.NewSource(sessionId))
		secondPlayerSession := MakePlayerSession(conn, secondPlayerPieceGenerator, session)
		session.SecondPlayerSession = secondPlayerSession
		session.FirstPlayerSession.EnemySession = secondPlayerSession
		session.SecondPlayerSession.EnemySession = session.FirstPlayerSession
		session.Started = true
		session.RunSession()
		log.Printf("Session %d started", sessionId)
	}

}
