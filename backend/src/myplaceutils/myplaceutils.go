package myplaceutils

import (
	//"strings"
	"fmt"
	"net"
	"reflect"
	"time"
  "log"
)


//Dessa loggers är till för att kunna anropas från alla programfiler.
var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

type User struct {
	Uname      string
	Pass       string
	Rooms      []*Room
	ActiveConn net.Conn
}

type Room struct {
	Name     string
	NoPeople int
	Users    []*User
	Messages []Message // Den kanske ska innehålla pekare till meddelanden?
}

type Message struct {
	Time  time.Time
	Uname string
	Body  string
	ID    string
}

//User method for binding the current connection to the user
func (u *User) BindConnection(c net.Conn) bool {
	u.ActiveConn = c
	return true
}

//Method for room to add a new message
//func (r Room)NewMessage(u User, msgbody string)

//Room method to add a user to the room
func (r *Room) AddUser(u *User) {
	r.Users = append(r.Users, u)
	r.NoPeople++
	u.JoinRoom(r)
}

//User method to add a room to the list of room that the user is part of
func (u *User) JoinRoom(r *Room) {
	u.Rooms = append(u.Rooms, r)

}

// Removes the user from the room
func (r *Room) RemoveUser(u *User) {
	for i, elem := range r.Users {
		if reflect.DeepEqual(elem, u) {
			r.Users = r.Users[:i+copy(r.Users[i:], r.Users[i+1:])]
		}
	}
	r.NoPeople--

}

// Removes the room from the user
func (u *User) RemoveRoom(r *Room) {
	for i, elem := range u.Rooms {
		if reflect.DeepEqual(elem, r) {
			u.Rooms = u.Rooms[:i+copy(u.Rooms[i:], u.Rooms[i+1:])]
		}

	}

}

func CreateUser(uname string, pass string, c net.Conn) User {
	u := User{}
	u.Uname = uname
	u.Pass = pass
	u.Rooms = []*Room{}
	u.ActiveConn = c
	return u
}

//Purpose: returns an array of the names of the rooms the user is in
//Use:    When the client software wants to list rooms, passing a name as an argument for joining a room, etc.
//Tested: NO
func (u User) showRooms() []string {
	var roomNames []string
	for _, r := range u.Rooms {
		roomNames = append(roomNames, r.Name)
	}
	return roomNames
}

//Purpose: Creating a new room
//Use: To create a new chat room
//Tested: No
func CreateRoom(name string) Room {
	return Room{name, 0, []*User{}, []Message{}}
}

//Purpose: returns an array of the names of the users in the room
//Use: when the client or server wishes to know what users are in the room
//Tested: NO
func ShowUsers(r Room) []string {
	var users []string
	for _, u := range r.Users {
		fmt.Printf("%v\n", u.Uname)
		users = append(users, u.Uname)
	}
	return users
}

func getUser(name string) {

}

func getRoom(id string) {

}

func getOlderMessages() {

}

func getNewerMessages() {

}

//Jävligt snyggt
func PrintTitle() {
	fmt.Printf("                    ____  __            \n")
	fmt.Printf("   ____ ___  __  __/ __ \\/ /___ _________ \n")
	fmt.Printf("  / __ `__ \\/ / / / /_/ / / __ `/ ___/ _ \\\n")
	fmt.Printf(" / / / / / / /_/ / ____/ / /_/ / /__/  __/\n")
	fmt.Printf("/_/ /_/ /_/\\__, /_/   /_/\\__,_/\\___/\\___/ \n")
	fmt.Printf("          /____/                          \n")

}
