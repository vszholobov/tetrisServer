package field

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
	"unicode/utf8"
)

type GameSession struct {
	firstPlayerSession  *PlayerSession
	secondPlayerSession *PlayerSession
}

type PlayerSession struct {
	playerField        *Field
	conn               *websocket.Conn
	playerInputChannel chan rune
	isEnded            bool
}

func MakePlayerSession(playerField *Field, conn *websocket.Conn) PlayerSession {
	return PlayerSession{playerField: playerField, conn: conn, playerInputChannel: make(chan rune), isEnded: false}
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
				// TODO: Signal "Game Over"
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
		// TODO: send field
		playerSession.conn.WriteMessage(websocket.TextMessage, []byte("1"+gameField.String()))
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
