package websockets

import (
	"bytes"
	"encoding/json"
	"net/http"
	"snap_chat_server/config"
	"snap_chat_server/logger"
	"snap_chat_server/services"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	auth *services.AuthSession

	id uuid.UUID

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(db *gorm.DB) {
	defer func() {
		logger.AppLog.Debugf("[client : %s] Remove Client.", c.id)
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg ClientMessage
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.AppLog.Errorf(err, "[client(%s) : %s] Connection closed", c.auth.Username, c.id)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		if err = json.Unmarshal(message, &msg); err != nil {
			logger.AppLog.Errorf(err, "[client(%s) : %s] Cannot unpacked message", c.auth.Username, c.id)
			break
		}
		logger.AppLog.Debugf("[client(%s) : %s] Send message : %s", c.auth.Username, c.id, msg.Message)

		if msg.MessageType == MT_PrivateChat {

			privateChat := services.NewPrivateChat(services.PrivateChatSender{
				Name:     c.auth.Name,
				Username: c.auth.Username,
			}, msg.Target[0])

			room, err := privateChat.FindOrCreate(db, *c.auth)

			if err != nil {

				errorMessage := ClientMessage{
					IsError: true,
					Sender: Sender{
						Name:     c.auth.Name,
						Username: c.auth.Username,
					},
					Target:  []string{c.auth.Username},
					Message: []byte(err.Error()),
				}

				byteMessage, _ := json.Marshal(errorMessage)

				c.send <- byteMessage

			} else {

				c.hub.chatIn <- ClientMessage{
					RoomUID:  room.RoomUID,
					RoomName: room.Name,
					Target:   msg.Target,
					Sender: Sender{
						Name:     c.auth.Name,
						Username: c.auth.Username,
					},
					Message: msg.Message,
				}

			}
		} else if msg.MessageType == MT_GroupChat {

			groupChat := services.NewGroupChat(services.GroupChatSender{
				Name:     c.auth.Name,
				Username: c.auth.Username,
			}, msg.Target[0])

			room, targetAudiences, err := groupChat.Find(db, *c.auth)

			if err != nil {

				errorMessage := ClientMessage{
					IsError: true,
					Sender: Sender{
						Name:     c.auth.Name,
						Username: c.auth.Username,
					},
					Target:  []string{c.auth.Username},
					Message: []byte(err.Error()),
				}

				byteMessage, _ := json.Marshal(errorMessage)

				c.send <- byteMessage

			} else {

				c.hub.chatIn <- ClientMessage{
					RoomUID:  room.RoomUID,
					RoomName: room.Name,
					Target:   targetAudiences,
					Sender: Sender{
						Name:     c.auth.Name,
						Username: c.auth.Username,
					},
					Message: msg.Message,
				}

			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(db *gorm.DB) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Debug("New Client Established.")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.AppLog.Fatal(err, "Failed to upgrade http serve into websocket")
		return
	}

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Replace(reqToken, "Bearer ", "", 1)
	reqToken = splitToken

	claims, err := services.VerifyToken(reqToken, config.Env.SecretKey)

	if err != nil {
		logger.AppLog.Fatal(err, "Verification token failed")
		return
	}

	username := claims["username"].(string)

	authSession, err := services.GetAuthUser(db, username, r)

	if err != nil {
		logger.AppLog.Fatal(err, "Get user session failed")
		return
	}

	clientId := uuid.New()

	logger.AppLog.Debugf("Register New Client with id : %s", clientId)
	client := &Client{
		id:   clientId,
		auth: authSession,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump(db)

	// Read messages
	logger.AppLog.Debugf("[client(%s) : %s] Message receiver ready", username, clientId)
	go client.readPump(db)
}
