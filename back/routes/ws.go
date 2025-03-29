package routes

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/ggmolly/galactf/orm"
	protobuf "github.com/ggmolly/galactf/proto"
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
	wsLock.Lock()
	defer wsLock.Unlock()

	// Register the client's ID with its connection
	Sockets[cid] = WsClient{
		UserID: user.ID,
		Conn:   c,
	}
}

func RemoveClient(cid uint64) {
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

func RevealAgent() {
	var challenges []orm.Challenge
	if err := orm.GormDB.
		Preload("Attachments").
		Order("reveal_at ASC").
		Where("reveal_at >= ?", time.Now().UTC()).
		Find(&challenges).Error; err != nil {
		log.Println("[!] failed to load challenges", err)
		return
	}

	for _, c := range challenges {
		now := time.Now().UTC()
		delay := c.RevealAt.Sub(now)
		log.Printf("[#] sleeping for %v before revealing challenge %d [%s]", delay, c.ID, c.Name)
		go func() {
			time.Sleep(delay+1*time.Second)
			attachments := make([]*protobuf.Attachment, len(c.Attachments))

			// Copy from ORM to protobuf
			for _, a := range c.Attachments {
				attachments = append(attachments, &protobuf.Attachment{
					Id:      a.ID,
					Type:    a.Type,
					Url:     a.URL,
					Filename: a.Title,
					Size:    a.Size,
				})
			}

			log.Printf("[#] revealing challenge %d [%s]", c.ID, c.Name)
			Broadcast(protobuf.WS_CHALLENGE_REVEAL, &protobuf.ChallengeReveal{
				Id:            c.ID,
				Name:          c.Name,
				Difficulty:    int32(c.Difficulty),
				Categories:    c.Categories,
				Attachments:   attachments,
				Description:   c.Description,
			})
		}()
	}
}