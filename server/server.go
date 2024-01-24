package server

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"tetrisServer/field"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// var Addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var Addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var sessions = make(map[int64]*field.GameSession)

//func Echo(w http.ResponseWriter, r *http.Request) {
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Print("upgrade:", err)
//		return
//	}
//	gameSession := field.MakeGameSession()
//	sessions[gameSession.GetSessionId()] = gameSession
//	//gameSession.RunSession()
//	//session := field.MakePlayerSession(conn)
//	//session.RunSession()
//}

type CreateSessionResponse struct {
	SessionId int64 `json:"sessionId"`
}

func GetSessionsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	gameSession := field.MakeGameSession()
	sessions[gameSession.GetSessionId()] = gameSession
	response := CreateSessionResponse{SessionId: gameSession.GetSessionId()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	//w.Write(big.NewInt(gameSession.GetSessionId()).Bytes())
	//homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func ConnectToSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId, _ := strconv.ParseInt(vars["sessionId"], 10, 64)
	session := sessions[sessionId]

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
		firstPlayerSession := field.MakePlayerSession(conn, firstPlayerPieceGenerator)
		session.FirstPlayerSession = firstPlayerSession
	} else {
		secondPlayerPieceGenerator := rand.New(rand.NewSource(sessionId))
		secondPlayerSession := field.MakePlayerSession(conn, secondPlayerPieceGenerator)
		session.SecondPlayerSession = secondPlayerSession
		session.FirstPlayerSession.EnemySession = secondPlayerSession
		session.SecondPlayerSession.EnemySession = session.FirstPlayerSession
		session.Started = true
		session.RunSession()
	}

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
