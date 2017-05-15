package myplaceutils

import (
	"net"
	"testing"
	"log"
//	"fmt"
//	"reflect"
)

func establishConnection() *net.TCPConn {
	tcpLAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1337")
	conn, err := net.DialTCP("tcp", nil, tcpLAddr) //localhost och port 1337

	if err != nil {
		log.Fatal(err)
	}
	return conn

}

func TestGetLatestMsg(t *testing.T){
	InitDBs()
	u1 := AddNewUser("ask", "embla")
	u2 := AddNewUser("adam", "eva")
	r1 := AddNewRoom("livingroom")
	u1.JoinRoom(r1)
	u2.JoinRoom(r1)

	returnLatestMsg, _ := GetLatestMsg(r1)
	
	if returnLatestMsg != nil {
		t.Error("The room shouldn't have any messages yet")
	}

	AddNewMessage(u1.UName,r1,"Hello, can you hear me?")
	AddNewMessage(u2.UName,r1,"Yes I Can")
	
	returnLatestMsg, _ = GetLatestMsg(r1)

	if returnLatestMsg.ID != 1 {
		t.Error("Didn't get latest message")
	}

	AddNewMessage(u1.UName,r1,"How are you?")
	AddNewMessage(u2.UName,r1,"Im fine thank you")
	AddNewMessage(u1.UName,r1,"Nice, have to go, bb!")

	returnLatestMsg, _ = GetLatestMsg(r1)
	
	if returnLatestMsg.ID != 4 {
		t.Error("Didn't get the latest message")

	}

	
}

/*
func createStuff(conn *net.TCPConn) (*User, *Room) {
	workingRoom := CreateRoom("Room 213")
	workingUser := CreateUser("Alex", "1337", conn)
	return workingUser, workingRoom
}

func TestJoinRoom(t *testing.T) {
	conn := establishConnection()
	workingUser, workingRoom := createStuff(conn)

	if len(workingUser.Rooms) != 0 {
		t.Error("The user isn't empty")
	}

	workingUser.JoinRoom(workingRoom)

	if len(workingUser.Rooms) != 1 {
		t.Error("Failed to update the user, user failed to join  213")
	}

}

func TestAddUser(t *testing.T) {
	conn := establishConnection()
	user1, workingRoom := createStuff(conn)

	if len(workingRoom.Users) != 0 && workingRoom.NoPeople != 0 {
		t.Error("Working room isn't empty from the start")
	}

	workingRoom.AddUser(user1)

	if len(workingRoom.Users) != 1 {
		t.Error("Failed to update the room, Room failed to get the user Alex ")
	}

	user2 := CreateUser("Erik", "1821", conn)
	workingRoom.AddUser(user2)

	if len(workingRoom.Users) != 2 {
		t.Error("Failed to update the room with 2 users ")
	}

}

func TestRemoveUser(t *testing.T) {
	conn := establishConnection()
	room := CreateRoom("Room 213")
	user0 := CreateUser("user0", "polis", conn)
	user1 := CreateUser("user1", "skurk", conn)
	user2 := CreateUser("user2", "inbrottstjuv", conn)

	room.AddUser(user0)
	room.AddUser(user1)
	room.AddUser(user2)

	if room.Users[2].Uname != "user2" {
		t.Error("The last elmenent isn't equal to user2")
	}

	room.RemoveUser(user2)

	// Kollar om user2 finns kvar i room.Users
	for _, elem := range room.Users {
		if reflect.DeepEqual(elem, user2) {
			t.Error("Didn't succed to remove the user")
		}
	}

	room.RemoveUser(user0)
	// Kollar om om user 0 finns kvar i rummet
	for _, elem := range room.Users {
		if reflect.DeepEqual(elem, user0) {
			t.Error("Didn't succed to remove the user")
		}
	}

	if room.NoPeople != 1 {
		t.Error("Number of users aren't equal to 1")
	}

}

func TestLeaveRoom(t *testing.T) {
	conn := establishConnection()
	user := CreateUser("MainUser", "1337", conn)
	room0 := CreateRoom("Room 0")
	room1 := CreateRoom("Room 1")
	room2 := CreateRoom("Room 2")

	user.JoinRoom(room0)
	user.JoinRoom(room1)
	user.JoinRoom(room2)

	if len(user.Rooms) != 3 {
		t.Error("Failed to update the user, user failed to join  213")
	}

	user.RemoveRoom(room0)

	for _, elem := range user.Rooms {
		if reflect.DeepEqual(elem, room0) {
			t.Error("Didn't succed to remove the room from the user")
		}
	}

	user.RemoveRoom(room2)

	for _, elem := range user.Rooms {
		if reflect.DeepEqual(elem, room2) {
			t.Error("Didn't succed to remove the room from the user")
		}
	}
}



func TestGetUser(t *testing.T){
	conn := establishConnection()
	usr0 := CreateUser("ussr0", "pass", conn)

	usr1 := GetUser("ussr0")

	usr1.Uname = "fisk" // kollar så att namnet ändras för alla instanser av usr
	
	if !(reflect.DeepEqual(usr0,usr1)){
		t.Error("User not in global User list")
	}

	usr2 := GetUser("MainUser")

	if (reflect.DeepEqual(usr1,usr2)){
		t.Error("Different users are the same")
	}
}

func TestShowUsers(t *testing.T){

	usr0 := GetUser("MainUser")

	rooms := usr0.ShowRooms()

	for x,y := range rooms{
		fmt.Printf("%d,%s",x,y)
	}
}

func TestPostMsg(t *testing.T){
	conn := establishConnection()
	usr0 := CreateUser("ussr1", "pass", conn)
	room0 := CreateRoom("Rum0")

	usr0.PostMsg(room0,"hej detta är ett meddelande")

	msg := room0.Messages[0]
	
	fmt.Printf("%v",msg)
}

*/

