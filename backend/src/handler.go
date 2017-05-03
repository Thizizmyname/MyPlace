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

var Users []User
var Rooms []Room


func handler(connection net.Conn, gochan chan string){
	Users,Rooms,_ = loadDBs()
}

func getUser(username string) User{

	for _,x := range User{
		if x.Uname == username{
			return x
		}
	}
	return nil
}

func getRoom(roomname string) Room{

	for _,x := range Rooms{
		if x.Name == roomname{
			return x
		}
	}
	return nil
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
	// ska användarna eller rummen ha channel i sig?
	//ME

	r := getRoom(roomID)
	m := getMsg(r, msgID)

	return []string{"he","he","he"}
}

func getNewerMessages(roomID string, msgID string) []string{

	r := getRoom(roomID)
	
	//
	//ME
	// vet inte riktigt hur denna ska funka.
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
	r := CreateRoom(roomName)
	u := getUser(username)
	r.addUser(u) 
	return r
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



