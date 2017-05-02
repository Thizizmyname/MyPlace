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
/*
00.signup
args: uname, pass
response: -
note: error if uname is taken/ pass to short/ illegal characters/ ...
side-effect: updates users_db
*/
func signup(username string, password string) {

}

/*
01.signin
args: uname, pass
response: -
note: error if uname not in use/ incorrect pass
side-effect: new messages from all rooms that the client has joined start being pushed to client
*/
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



