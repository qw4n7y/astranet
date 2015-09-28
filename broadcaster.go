package astranet

import (
  "net/http"
  "fmt"
)

type Broadcaster struct {
  connections   map[*Connection]bool
}

func (b *Broadcaster) Setup() {
  go b.loopBroadcasting()
}

func (b *Broadcaster) loopBroadcasting() {
  for {
    for connection, _ := range b.connections {
      select {
      case data, ok := (<-connection.read):
        if (!ok) {
          b.Remove(connection)
        }

        b.Broadcast(data)

      default:
        // nothing there. just looping
      }
    }
  }
}

//  Broadcasts data to every connection
//
func (b *Broadcaster) Broadcast(data []byte) {
  for connection, _ := range b.connections {
    select {
      case connection.write <- data:
        // sent data to write pipe

      default:
        // so agressive?
        b.Remove(connection)
    }
  }
}

//  Adds new connection
//
func (b *Broadcaster) Add(connection *Connection) {
  b.connections[connection] = true
  connection.Setup()

  fmt.Printf("%s was added\n", connection)
}

//  Removes new connection
//
func (b *Broadcaster) Remove(connection *Connection) {
  delete(b.connections, connection)
  connection.Close()

  fmt.Printf("%s was removed\n", connection)
}

// Handles websocket requests from the peer
//
func serveWebsockets(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return
  }

  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    fmt.Println(err)
    return
  }

  connection := &Connection{ ws: ws }

  iBroadcaster.Add(connection)
}