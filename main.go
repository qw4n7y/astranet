package astranet

import (
  "net/http"
  "fmt"
  )

//  Global variables
//
var iBroadcaster *Broadcaster

//  The magic starts here
//
func main() {
  fmt.Println("Hello, Universe!")

  iBroadcaster = new(Broadcaster)
  iBroadcaster.Setup()
  
  http.Handle("/", http.FileServer(http.Dir("./public")))
  http.HandleFunc("/ws", serveWebsockets)
  
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("Error: " + err.Error())
  }
}