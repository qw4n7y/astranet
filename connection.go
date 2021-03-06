package astranet

import (
  "github.com/gorilla/websocket"
  "code.google.com/p/go-uuid/uuid"
  "time"
  "fmt"
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

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

type Connection struct {
  ws *websocket.Conn
  uuid string

  write chan []byte
  read chan []byte
}

// To String
//
func (c *Connection) String() string {
  return fmt.Sprintf("<Connection %s>", c.uuid)
}

//  Close
//
func (c *Connection) Close() {
  close(c.read)
  close(c.write)

  c.ws.Close()
}

//  Write to socket
//
func (c *Connection) writeMessage(messageType int, payload []byte) error {
  c.ws.SetWriteDeadline(time.Now().Add(writeWait))
  return c.ws.WriteMessage(messageType, payload)
}

//  Writing loop
//
func (c *Connection) loopWriting() {
  ticker := time.NewTicker(pingPeriod)

  defer func() {
    ticker.Stop()
    c.Close()
  }()
  
  for {
    select {
    case message, ok := <- c.write:
      if !ok {
        c.writeMessage(websocket.CloseMessage, []byte{})
        return
      }
      if err := c.writeMessage(websocket.TextMessage, message); err != nil {
        return
      }
    case <-ticker.C:
      if err := c.writeMessage(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
  }
}

//  Reading loop
//
func (c *Connection) loopReading() {
  defer func() {
    c.Close()
  }()

  c.ws.SetReadLimit(maxMessageSize)
  c.ws.SetReadDeadline(time.Now().Add(pongWait))
  c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

  for {
    _, message, err := c.ws.ReadMessage()
    if err != nil {
      break
    }
    c.read <- message 
  }
}

//  Setting up
//
func (c *Connection) Setup() {
  c.read = make(chan []byte)
  c.write = make(chan []byte)

  c.uuid = uuid.New()

  go c.loopReading()
  go c.loopWriting()
}