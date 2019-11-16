package server
import (
	"ChatServer/src/server/components"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Handler struct {}

var reqHandler = &Handler{}
var chatConnectionManger = &components.ConnectionMgrImpl{ActiveConnections: make(map[*websocket.Conn] *components.ConnectionImpl)}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if err != nil {
		panic("Could not upgrade http to ws connection")
	}
	chatConnectionManger.RegisterConnection(conn)
	log.Println("Received a ws connection req")
}

func StartChatServer() {
	http.HandleFunc("/ws", reqHandler.ServeHTTP)
	if err:= http.ListenAndServe(":8080", nil) ; err != nil {
		panic("Cannot start chat server ...")
	}
}
