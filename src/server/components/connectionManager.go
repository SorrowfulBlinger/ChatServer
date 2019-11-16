package components

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
)

type ConnectionManager interface {
	setup()
	registerConnection(conn *websocket.Conn)
	UnregisterConnection()
	Broadcast([]byte)
}

type ConnectionMgrImpl struct {
	ActiveConnections map[*websocket.Conn] *ConnectionImpl
}

func (connMgr *ConnectionMgrImpl) RegisterConnection(conn *websocket.Conn) {
	id, err := uuid.NewV1()
	if err != nil {
		log.Fatal("Cannot generate UUID")
	}
	connImpl := &ConnectionImpl{
		id:              id.String(),
		connMgr:         connMgr,
		wsConnection:    conn,
		toClientChannel: make(chan []byte),
	}
	connImpl.Setup()
	connMgr.ActiveConnections[conn] = connImpl
	log.Printf("successfully registered %v", id.String())
}

func (connMgr *ConnectionMgrImpl) UnregisterConnection(conn *websocket.Conn) {
	if _, ok := connMgr.ActiveConnections[conn]; ok {
		log.Printf("successfully unregistered %v", connMgr.ActiveConnections[conn].id)
		delete(connMgr.ActiveConnections, conn)
	} else {
		log.Printf("already unregistered")
	}
}

func (*ConnectionMgrImpl) Setup() {
}

func (connMgr *ConnectionMgrImpl) Broadcast(msg []byte, sender *ConnectionImpl) {
	// index, val returned
	for _, connection := range connMgr.ActiveConnections {
		if connection != sender {
			connection.Send(msg)
		}
	}
}

