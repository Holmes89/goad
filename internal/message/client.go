package message

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client interface{
	Send(message *Message)
	Run()
	Close()
}

type grpcClient struct {
	hub Hub
	messenger Messenger_SendMessageServer
	wg sync.WaitGroup
}

type wsClient struct {
	hub Hub
	conn *websocket.Conn
	send chan *Message
}


func NewGRPCClient(hub Hub, messenger Messenger_SendMessageServer) Client {
	c := &grpcClient{
		hub: hub,
		messenger: messenger,
	}
	logrus.Info("grpc client created")
	return c
}

func NewWebSocketClient(hub Hub, w http.ResponseWriter, r *http.Request) Client {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithField("err", err).Error("unable to register ws client")
		return nil
	}
	s := make(chan *Message)
	c := &wsClient{
		hub: hub,
		conn: conn,
		send: s,
	}

	logrus.Info("ws client created")
	return c
}

func (c *grpcClient) Send(message *Message) {
	c.messenger.Send(message)
}

func (c *grpcClient) Close() {
	c.hub.Unregister(c)
}

func (c *grpcClient) Run() {
	for {
		resp, err := c.messenger.Recv()
		if err == io.EOF {
			logrus.Error("eof")
			c.Close()
			break
		}
		if err != nil {
			handleErrorMessage(err)
			c.Close()
			break
		}

		//TODO Add uuid, timestamp
		message := &Message{
			From: resp.From,
			Body: resp.Body,
		}
		c.hub.Broadcast(message)
	}
}

func (c *wsClient) Send(message *Message) {
	c.send <- message
}

func (c *wsClient) Close() {
	close(c.send)
	c.hub.Unregister(c)
}

func (c *wsClient) Run() {
	go c.readPump()
	go c.writePump()
}

func (c *wsClient) readPump(){
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, resp, err := c.conn.ReadMessage()
		if err != nil {
			logrus.WithField("error", err.Error()).Error("websocket error")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithField("error", err.Error()).Error("websocket error")
			}
			break
		}

		//TODO Add uuid, timestamp

		//message := &Message{
		//	From: "Web User",
		//	Body: string(resp),
		//}

		message := &Message{}
		if err := json.Unmarshal(resp, message); err != nil {
			logrus.Error("unable to parse")
		}

		c.hub.Broadcast(message)
	}
}

func (c *wsClient) writePump(){

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			logrus.Info("sending message from ws")
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
			msg := fmt.Sprintf("%s: %s", message.From, message.Body)
			w.Write([]byte(msg))
			w.Write([]byte{'\n'})

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


func handleErrorMessage(err error) {
	status, _ := status.FromError(err)
	if codes.Canceled == status.Code() {
		logrus.Info("connection closed")
	} else {
		logrus.WithField("error", status.Message()).Error("received error from client")
	}
}

