package myplaceutils

import (
    //"strings"
    "fmt"
    "net"
    "time"
)


type User struct {
	Uname string
	Pass string
	Rooms []Room
  ActiveConn net.Conn
}

type Room struct {
	Name string
	NoPeople int
  Users []User
	Messages []Message
}

type Message struct {
	Time time.Time
	Uname string
	Body string
  ID string
}

//User method for binding the current connection to the user
func (u *User)NewConnection(c net.Conn) bool {
  u.ActiveConn = c
  return true
}

//Method for room to add a new message
//func (r Room)NewMessage(u User, msgbody string)


//Room method to add a user to the room
func (r Room)AddUser(u User){
  r.Users = append(r.Users, u)
  u.JoinRoom(r)
}


//User method to add a room to the list of room that the user is part of
func (u *User)JoinRoom(r Room){
  u.Rooms = append(u.Rooms, r)

}

func CreateUser(uname string, pass string, c net.Conn) User{
  u := User{}
  u.Uname = uname
  u.Pass = pass
  u.Rooms = []Room{}
  u.ActiveConn = c
  return u
}


func CreateRoom(name string) Room{
  return Room{name, 0, []User{}, []Message{}}
}

func (r Room)ShowUsers(){
  for i,u := range r.Users {
    fmt.Printf("User %d: %v\n", i, u.Uname)
  }
}
