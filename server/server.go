package server

import (
	"flag"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"tetrisServer/field"
	"time"
)

var Addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func Echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	gameField := field.MakeDefaultField()
	session := field.MakePlayerSession(&gameField, conn)
	session.RunSession()

	//keyboardSendChannel := make(chan rune)
	//go func(gameField *field.Field, keyboardInputChannel chan rune) {
	//	for {
	//		inputControl(keyboardInputChannel, gameField)
	//
	//		if !gameField.CurrentPiece.MoveDown() {
	//			gameField.Val.Or(gameField.Val, gameField.CurrentPiece.GetVal())
	//			gameField.SelectNextPiece()
	//			if !gameField.CurrentPiece.CanMoveDown() {
	//				// TODO: Player lost
	//				field.CallClear()
	//				fmt.Println("Game over. Stats:")
	//				fmt.Printf("Score: %d | Speed: %d | Lines Cleand: %d\n", *gameField.Score, gameField.GetSpeed(), *gameField.CleanCount)
	//				break
	//			}
	//			gameField.Clean()
	//		}
	//	}
	//}(&extField, keyboardSendChannel)
	//go func(conn *websocket.Conn) {
	//	for {
	//		_, message, err := conn.ReadMessage()
	//		if err != nil {
	//			//log.Println("read:", err)
	//			break
	//		}
	//		decodeRune, _ := utf8.DecodeRune(message)
	//		keyboardSendChannel <- decodeRune
	//		log.Printf("recv: %s", message)
	//		//err = conn.WriteMessage(mt, message)
	//		//if err != nil {
	//		//	log.Println("write:", err)
	//		//	break
	//		//}
	//	}
	//}(conn)
}

func inputControl(
	keyboardInputChannel chan rune,
	gameField *field.Field,
) {
	timeout := time.After(time.Second / 4 / time.Duration(gameField.GetSpeed()))
	for {
		field.PrintField(gameField)
		select {
		case moveType := <-keyboardInputChannel:
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
				gameField.CurrentPiece.Rotate(field.Left)
			case 101:
				// e
				gameField.CurrentPiece.Rotate(field.Right)
			}
		case <-timeout:
			return
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
