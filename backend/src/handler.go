package main

import (
    "net"
//    "fmt"
  //"time"
//  "bufio"
//  "log"
//  "reflect"
  //"data" //importera denna när vi ska implementera databasen
	//"myplaceutils"
)

//var connections []net.Conn


func handler(connection net.Conn, gochan chan string){
	
}

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
  return ""
}

func getRooms(username string /*Kanske byta mot user-type*/) []string{
	
	return []string{"h","a","ha"}
}

func getOlderMessages(roomID string, msgID string) []string /*kanske någon annan returntype*/ {
	// ska användarna eller rummen ha channel i sig?
	//ME



	return []string{"he","he","he"}
}

func getNewerMessages(roomID string, msgID string) []string{

	
	//
	//ME
	// vet inte riktigt hur denna ska funka.
  return []string{"ho","ho","ho"}
}

func joinRoom(roomId string, username string){
// Exempel på hur det kan se ut.	
//	workingUser := myplaceutils.getUser(username)
//	workingRoom := getRoom(roomId)
//	workingUser.JoinRoom(workingRoom)
	
}

func leaveRoom(username string, roomId string) {
/*
* workingUser := myplaceutils.getUser(username)
* workingRoom := getRoom(roomId)
* workingUser.LeaveRoom(workingRoom)
*/
}

func createRoom(username string, roomName string) {	
	//r := CreateRoom(roomName)
	//u := getUser(username)
	//r.addUser(u) 
	//return r
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



