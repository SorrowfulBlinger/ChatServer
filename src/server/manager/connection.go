package manager

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Connection interface {
	Setup()
	Send([]byte)bool
	Close()
}

type ConnectionImpl struct {
	id string
	connMgr *ConnectionMgrImpl
	wsConnection *websocket.Conn
	toClientChannel chan []byte
	closeLock sync.Mutex
}

func (conn *ConnectionImpl) Setup () {
	go conn.listenFromClient()
	go conn.pushToClient()
}

func (conn *ConnectionImpl) Close () {
	conn.closeLock.Lock()
	defer conn.closeLock.Unlock()
	if conn.wsConnection != nil {
		conn.wsConnection.Close()
		close(conn.toClientChannel)
		conn.connMgr.UnregisterConnection(conn.wsConnection)
		conn.wsConnection = nil
		log.Printf("closing connection %v", conn.id)
	}
}

func (conn *ConnectionImpl) listenFromClient() {
	defer conn.Close()
	label: for {
		_, msgReceived, err := conn.wsConnection.ReadMessage()
		if err != nil {
			log.Println(fmt.Errorf("cannot read from connection %v", conn.id))
			break label
		}
		conn.connMgr.Broadcast(msgReceived, conn)
	}
}

func (conn *ConnectionImpl) pushToClient() bool {
	defer conn.Close()
	label:
		for msg:= range conn.toClientChannel {
			if conn.wsConnection != nil {
				if err := conn.wsConnection.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println(fmt.Errorf("cannot write to connection %v", conn.id))
					break label
				}
			}
		}
	return true
}

func (conn *ConnectionImpl) Send (message []byte) bool {
	conn.toClientChannel <- message
	return true
}
