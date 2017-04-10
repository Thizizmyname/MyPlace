package main

import (
    "net"
    "fmt"
)

func main() {
  fmt.Printf("Starting server...\n")
  listener, err := net.Listen("tcp",":1337")
  if err != nil {
    fmt.Printf("Error starting server: %v\n", err)
  }
  for {
    connection, errs := listener.Accept()
    if errs != nil {
      fmt.Printf("Connection Accept error: %v\n",errs)
    } else {
      fmt.Printf("Connection established: %v\n", connection)
      go handleConnection()
    }
  }
}


func handleConnection() {
  fmt.Printf("Handle connection whattup\n")
}

