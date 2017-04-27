package main

import (
    "net"
//    "fmt"
  //"time"
//  "bufio"
//  "log"
//  "reflect"
  //"data" //importera denna när vi ska implementera databasen
//  "myplaceutils"
)



func handler(connection net.Conn, gochan chan string){
   //
}


func signup(username string, password string) {

}


func signin(username string, password string) /*TOKEN*/ string{
  //throw error
  return ""
}

func getRooms(username string /*Kanske byta mot user-type*/) []string{
  return []string{"h","a","ha"}
}

func getOlderMessages(roomID string, msgID string) []string /*kanske någon annan returntype*/ {
  //
  return []string{"he","he","he"}
}

func getNewerMessages(roomID string, msgID string) []string{
  //
  return []string{"ho","ho","ho"}
}

func joinRoom(roomId string, username string) {
// Exempel på hur det kan se ut.	
//	workingUser = getUser(username)
//	workingRoom = getRoom(roomId)
//	workingUser.JoinRoom(workingRoom)
	
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



