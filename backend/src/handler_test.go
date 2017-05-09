package main

import (

	"testing"
	"net"
	"log"
)




func establishConnection() *net.TCPConn {
	tcpLAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1337")
	conn, err := net.DialTCP("tcp", nil, tcpLAddr) //localhost och port 1337

	if err != nil {
		log.Fatal(err)
	}
	return conn

}

/*
00.signup
args: uname, pass
response: -
note: error if uname is taken/ pass to short/ illegal characters/ ...
side-effect: updates users_db
*/

func TestSignUp(t *testing.T){
	conn := establishConnection()
	SignUp("Alex","hej123",conn)

	if SignUp("Erik","1337",conn) {
		
	}

	if SignUp("Alex","lol123",conn) {
		t.Error("Username allready taken\n")
	}


	
	
}


func TestJoinRoom(t *testing.T){
	
	

}
