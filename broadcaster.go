package main

type Broadcaster struct {
  connections   map[*Connection]bool

  add           chan *Connection
}

func (b *Broadcaster) Setup() {
  go b.loopAdding()
  go b.loopBroadcasting()
}

func (b *Broadcaster) loopAdding() {
  for {
    select {
      case connection, ok := (<-add):
        if (ok) {
          b.addConnection(connection)
        }
    }
  }
}

func (b *Broadcaster) loopBroadcasting() {
  for {
    for connection, _ := range b.connections {
      select {
      case data, ok := (<-connection.read):
        if (!ok) {
          b.removeConnection(connection)
        }

        b.Broadcast(data)

      default:
        // nothing there. just looping
      }
    }
  }
}

func (b *Broadcaster) Broadcast(data) {
  for connection, _ := range b.connections {
    select {
      case connection.write <- data:
        // sent data to write pipe

      default:
        // so agressive?
        b.removeConnection(connection)
    }
  }
}

func (b *Broadcaster) addConnection(connection) {
  b.connections[connection] = true
  connection.Setup()
}

func (b *Broadcaster) removeConnection(connection) {
  delete(b.connections, connection)
  connection.Close()
}

//  Global variables
//
var iBroadcaster := new(Broadcaster)
iBroadcaster.Setup()

// Handles websocket requests from the peer
//
func serveWebsockets(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return
  }

  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  connection := &Connection{ ws: ws }

  iBroadcaster.add <- connection
}