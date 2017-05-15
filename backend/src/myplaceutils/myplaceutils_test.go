package myplaceutils

import (
	"net"
	"testing"
	"log"
	//"fmt"
	//"reflect"
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
	workingUser := CreateUser("Alex", "1337")
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

	user2 := CreateUser("Erik", "1821")
	workingRoom.AddUser(user2)

	if len(workingRoom.Users) != 2 {
		t.Error("Failed to update the room with 2 users ")
	}

}

func TestRemoveUser(t *testing.T) {
	conn := establishConnection()
	room := CreateRoom("Room 213")
	user0 := CreateUser("user0", "polis")
	user1 := CreateUser("user1", "skurk")
	user2 := CreateUser("user2", "inbrottstjuv")

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
	user := CreateUser("MainUser", "1337")
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

*/

func TestGetUser(t *testing.T){
	InitDBs()
	
	//conn := estalishConnection()
	usr0 := AddNewUser("ussr0", "pass")

	usr1 := GetUser("ussr0")

	usr1.UName = "fisk" // kollar så att namnet ändras för alla instanser av usr
	
	if usr0 != usr1 {
		t.Error("User not in global User list")
	}

	usr0 = GetUser("Not a valid username")

	if usr0 != nil{
		t.Error("This should be nil")
	}

}

func TestShowRooms(t *testing.T){
	InitDBs()

	usr0 := AddNewUser("usr0", "pass0")
	RoomNames := []string{"Room0","Room1","Room2","Room3"}
	
	room0 := AddNewRoom(RoomNames[0])
	room1 := AddNewRoom(RoomNames[1])
	room2 := AddNewRoom(RoomNames[2])
	room3 := AddNewRoom(RoomNames[3])

	usr0.JoinRoom(room0)
	usr0.JoinRoom(room1)
	usr0.JoinRoom(room2)
	usr0.JoinRoom(room3)

	Rooms := usr0.ShowRooms()

	for x,y := range Rooms {
		if RoomNames[x] != y {
			t.Errorf("RoomNames are broken on index %v", x)
		}
	}

	usr1 := AddNewUser("usr1", "pass1")
	Rooms = usr1.ShowRooms()

	if len(Rooms) != 0{
		t.Errorf("length should be 0. Actual:%v",len(Rooms))
	}

	
	FakeUsr := CreateUser("fakeusr", "fakepass")

	Rooms = FakeUsr.ShowRooms()

	if Rooms != nil{
		t.Error("Should return nil")
	}
}

func TestShowUsers(t *testing.T){
	InitDBs() 

	usrnames := []string{"usr1","usr2","usr3","usr4"}
	
	room1 := AddNewRoom("Room1")

	usr1 := AddNewUser(usrnames[0], "pass1")
 	usr2 := AddNewUser(usrnames[1], "pass2")
	usr3 := AddNewUser(usrnames[2], "pass3")
	usr4 := AddNewUser(usrnames[3], "pass4")
	
	usr1.JoinRoom(room1)
	usr2.JoinRoom(room1)
	usr3.JoinRoom(room1)
	usr4.JoinRoom(room1)
	
	names := ShowUsers(room1)

	for x,y := range names {
		if usrnames[x] != y {
			t.Errorf("UserNames are broken on index %v", x)
		}
	}

	
	room2 := AddNewRoom("Room2")
	names = ShowUsers(room2)

	if len(names) != 0{
		t.Errorf("length should be 0. Actual:%v",len(names))
	}

	RealFakeRoom := CreateRoom("fakeroom", -1)

	names = ShowUsers(RealFakeRoom)

	if names != nil{
		t.Error("Should return nil")
	}
	
}
