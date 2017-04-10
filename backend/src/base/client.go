package main

import (
    "net"
    "fmt"
    "bufio"
    )


func main() {
  fmt.Printf("Connecting to server..\n")
  //Skapar två nya variabler, conn(ection) och err(or)
  // err är nil om allt gick bra, annars inte nil
  conn, err := net.Dial("tcp","127.0.0.1:1337") //localhost och port 1337
    if err != nil { //om något gick fel 
      //Handle errors
      fmt.Printf("Error, %v", err)
    } else { //annars fortsätt
      fmt.Printf("Connection established to the server..\n")
      fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
        status, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
          fmt.Print("ERROR ERROR ERROR, CALL GO-CHAN")
        } else {
          fmt.Printf("status: %v",status)
        }
    }
}
