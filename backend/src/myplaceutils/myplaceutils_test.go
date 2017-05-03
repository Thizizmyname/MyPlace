package myplaceutils

import(
	"net"
	"testing"
	"log"
)

func establishConnection() *net.TCPConn{
	tcpLAddr,_ := net.ResolveTCPAddr("tcp","127.0.0.1:1337")
	conn, err := net.DialTCP("tcp",nil,tcpLAddr) //localhost och port 1337
	
	if(err != nil){
		log.Fatal(err)
	}
	return conn

}

func createStuff(conn *net.TCPConn) (User,Room){
	workingRoom := CreateRoom("Room 213")
	workingUser := CreateUser("Alex", "1337",conn)
	return workingUser, workingRoom
}

func TestJoinRoom(t *testing.T){

	conn := establishConnection()
	workingUser, workingRoom := createStuff(conn)
	workingUser.JoinRoom(&workingRoom)

	if len(workingUser.Rooms) != 1 {
		t.Error("Failed to update the user, user failed to join  213")
	}
	
}

func TestAddUser(t *testing.T){
	conn := establishConnection()
	user1, workingRoom := createStuff(conn)

	if len(workingRoom.Users) != 0 && workingRoom.NoPeople != 0 {
		t.Error("Working room isn't empty from the start")
	}

	workingRoom.AddUser(&user1)
	
	if len(workingRoom.Users) != 1 {
		t.Error("Failed to update the room, Room failed to get the user Alex ")
	}
	
	user2 := CreateUser("Erik", "1821", conn)
	workingRoom.AddUser(&user2)

	if len(workingRoom.Users) != 2 {
		t.Error("Failed to update the room with 2 users ")
	}

	
}

