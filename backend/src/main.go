package main

import (
    "net"
    "fmt"
  //"time"
  "bufio"
//  "log"
  "reflect"
  //"data" //importera denna när vi ska implementera databasen
  "myplaceutils"
)


var connections []net.Conn

func signup(username string, password string) {

}


func signin(username string, password string) /*TOKEN*/ string{
  //throw error
}

func getRooms(username string /*Kanske byta mot user-type*/) []string{
  
}

func getOlderMessages(roomID string, msgID string) []string /*kanske någon annan returntype*/ {
  //
}

func getNewerMessages(roomID string, msgID string) []string{
  //
}

func joinRoom(roomId string, username string) {
  //
}

func leaveRoom(username string, roomId string) {
  //
}

func createRoom(username string, roomName string) {
  //
}

func postMessage(username string, roomID string, text string){
  //
}

func messageRead(username string, roomID string, msgID string) {
  //
}



func deleteUser(username string) bool{
  return true
}



