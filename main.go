package main

import (
  "net/http"
  "fmt"
  )

func main() {
  fmt.Println("Hello, Universe!")
  
  http.Handle("/", http.FileServer(http.Dir("./public")))
  
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("Error: " + err.Error())
  }
}