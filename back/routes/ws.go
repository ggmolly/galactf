package routes

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/ggmolly/galactf/orm"
	"github.com/gofiber/contrib/websocket"
	"google.golang.org/protobuf/proto"
)

type WsClient struct {
	UserID uint64
	Conn   *websocket.Conn
}

var (
	wsLock  = sync.RWMutex{}
	Sockets = make(map[uint64]WsClient)
)

// A single user can be connected twice (or more) to the websocket
// this function derives his ID with the current time to a xxhash
func generateClientID(user *orm.User) uint64 {
	var buf bytes.Buffer
	t := time.Now().UnixNano()
	for i := 0; i < 8; i++ {
		buf.WriteByte(byte(t >> (i * 8)))
	}
	for i := 0; i < 8; i++ {
		buf.WriteByte(byte(user.ID >> (i * 8)))
	}
	cid := xxhash.Sum64(buf.Bytes())
	return cid
}

func RegisterClient(c *websocket.Conn, user *orm.User, cid uint64) {
	log.Printf("[ws] #%d connected, cid: %d\n", user.ID, cid)

	wsLock.Lock()
	defer wsLock.Unlock()

	// Register the client's ID with its connection
	Sockets[cid] = WsClient{
		UserID: user.ID,
		Conn:   c,
	}
}

func RemoveClient(cid uint64) {
	log.Printf("[ws] disconnected, cid: %d\n", cid)

	wsLock.Lock()
	defer wsLock.Unlock()
	delete(Sockets, cid)
}

// Broadcast to everyone
func Broadcast(eventId uint8, msg proto.Message) {
	if msg == nil {
		return
	}

	var buf bytes.Buffer
	buf.WriteByte(eventId)

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("[-] error marshalling event: %s", err.Error())
		return
	}

	buf.Write(msgBytes)

	wsLock.Lock()
	defer wsLock.Unlock()

	for _, c := range Sockets {
		c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	}
}

// Broadcast to everyone except the user that triggered the event
func BroadcastExcl(eventId uint8, msg proto.Message, user *orm.User) {
	if msg == nil {
		return
	}

	var buf bytes.Buffer
	buf.WriteByte(eventId)

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("[-] error marshalling event: %s", err.Error())
		return
	}

	buf.Write(msgBytes)

	wsLock.Lock()
	defer wsLock.Unlock()

	for _, c := range Sockets {
		if c.UserID == user.ID {
			continue
		}
		c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	}
}

// Broadcast to a specific user
func BroadcastTo(eventId uint8, msg proto.Message, user *orm.User) {
	if msg == nil {
		return
	}

	var buf bytes.Buffer
	buf.WriteByte(eventId)

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("[-] error marshalling event: %s", err.Error())
		return
	}

	buf.Write(msgBytes)

	wsLock.Lock()
	defer wsLock.Unlock()

	for _, c := range Sockets {
		if c.UserID == user.ID {
			c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
		}
	}
}

func WsHandler(c *websocket.Conn) {
	user := c.Locals("user").(*orm.User)
	cid := generateClientID(user)
	RegisterClient(c, user, cid)

	// Poll until a message is received (if we do, kill the client, we don't want him to talk to us)
	// or until the connection is closed
	for {
		c.ReadMessage() // blocking call here
		break
	}

	RemoveClient(cid)
}
