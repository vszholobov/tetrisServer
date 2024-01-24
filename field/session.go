package field

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

type GameSession struct {
	sessionId           int64
	FirstPlayerSession  *PlayerSession
	SecondPlayerSession *PlayerSession
	Started             bool
}

type PlayerSession struct {
	playerField        *Field
	conn               *websocket.Conn
	playerInputChannel chan rune
	isEnded            bool
	pieceGenerator     *rand.Rand
	EnemySession       *PlayerSession
	mu                 sync.Mutex
}

func MakeGameSession() *GameSession {
	sessionId := time.Now().Unix()
	return &GameSession{
		sessionId: sessionId,
	}
}

func MakePlayerSession(conn *websocket.Conn, pieceGenerator *rand.Rand) *PlayerSession {
	field := MakeDefaultField(pieceGenerator)
	session := PlayerSession{
		playerField:        &field,
		conn:               conn,
		playerInputChannel: make(chan rune),
		isEnded:            false,
		pieceGenerator:     pieceGenerator,
	}
	return &session
}

func (gameSession *GameSession) RunSession() {
	gameSession.FirstPlayerSession.RunSession()
	gameSession.SecondPlayerSession.RunSession()
}

func (gameSession *GameSession) GetSessionId() int64 {
	return gameSession.sessionId
}

func (playerSession *PlayerSession) SendMessage(message string) {
	playerSession.mu.Lock()
	defer playerSession.mu.Unlock()
	playerSession.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (playerSession *PlayerSession) RunSession() {
	go playerSession.processPlayerInput()
	go playerSession.processGameField()
}

func (playerSession *PlayerSession) processGameField() {
	gameField := playerSession.playerField
	for {
		playerSession.inputControl()

		if !gameField.CurrentPiece.MoveDown() {
			gameField.Val.Or(gameField.Val, gameField.CurrentPiece.GetVal())
			gameField.SelectNextPiece()
			if !gameField.CurrentPiece.CanMoveDown() {
				playerSession.isEnded = true
				playerSession.conn.WriteMessage(websocket.TextMessage, []byte("0"))
				playerSession.conn.Close()
				//endMessage := "Game over. Stats:"
				//fmt.Println(endMessage)
				//fmt.Printf("Score: %d | Speed: %d | Lines Cleand: %d\n", *gameField.Score, gameField.GetSpeed(), *gameField.CleanCount)
				break
			}
			gameField.Clean()
		}
	}
}

func (playerSession *PlayerSession) processPlayerInput() {
	for !playerSession.isEnded {
		// TODO: ticker чтобы не зависнуть когда сессия закончилась
		_, message, err := playerSession.conn.ReadMessage()
		if err != nil {
			//log.Println("read:", err)
			break
		}
		decodeRune, _ := utf8.DecodeRune(message)
		playerSession.playerInputChannel <- decodeRune
		log.Printf("recv: %s", message)
		//err = conn.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}
}

func (playerSession *PlayerSession) inputControl() {
	gameField := playerSession.playerField
	timeout := time.After(time.Second / 4 / time.Duration(gameField.GetSpeed()))
	for {
		//PrintField(gameField)
		// TODO: send enemy field to enemy session with id (self = 0, enemy = 1)
		message := fmt.Sprintf("%d %s %d %d %d %d", 1, gameField.String(), gameField.GetSpeed(), *gameField.Score, *gameField.CleanCount, gameField.NextPiece.pieceType)
		playerSession.SendMessage("0 " + message)
		playerSession.EnemySession.SendMessage("1 " + message)
		select {
		case moveType := <-playerSession.playerInputChannel:
			switch moveType {
			case 100:
				// d
				gameField.CurrentPiece.MoveLeft()
			case 97:
				// a
				gameField.CurrentPiece.MoveRight()
			case 115:
				// s
				gameField.CurrentPiece.MoveDown()
			case 113:
				// q
				gameField.CurrentPiece.Rotate(Left)
			case 101:
				// e
				gameField.CurrentPiece.Rotate(Right)
			}
		case <-timeout:
			return
		}
	}
}
