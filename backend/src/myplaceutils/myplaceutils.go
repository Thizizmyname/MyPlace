package myplaceutils

import (
	"fmt"
	"net"
	"time"
	//"reflect"
	"log"
	"container/list"
	"requests_responses"
)


//Dessa loggers är till för att kunna anropas från alla programfiler.
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	//connections []net.Conn
	Users UserDB
	Rooms RoomDB
)

type User struct {
	UName string
	Pass string
	Rooms *list.List //list of ints:roomID
}

type Room struct {
	ID int
	Name     string
	Users *list.List //list of strings:UName
	Messages map[int]Message //ID is key
	OutgoingChannels *list.List //[]chan requests_responses.Response
}

type Message struct {
	ID    int
	Time  time.Time
	Uname string
	Body  string
}

type HandlerArgs struct {
	IncomingRequest requests_responses.Request
	ResponseChannel chan requests_responses.Response
}

type UserDB map[string]User //UName is key
type RoomDB map[int]Room //ID is key

func InitDBs() {
	Users = make(map[string]User)
	Rooms = make(map[int]Room)
}

/*
Funkar inte för att data.go är trasig
// Loads DB to global variables
func Initialize(){
	Users,Rooms,err = LoadDbs()
	//ska hantera error på nåt sätt
}

//Store global variables to DB
func Terminate(){
	err := StoreDBs(Users,Rooms)
	//ska hantera error på nåt sätt
}
*/



//MÅSTE HA MUTEX LOCK I DENNA FUNKTIONEN
//MÅSTE KOLLA SÅ INTE newConnection REDAN FINNS I connections
func AddConnection(newConnection net.Conn) {
	//connections = append(connections, newConnection)
}

//Kolla genom arrayen om den finns innan den försöker ta bort den
func RemoveConnection(connection net.Conn) bool{
	return true
}

//User method for binding the current connection to the user
func (u *User) BindConnection(c net.Conn) bool {
	//u.ActiveConn = c
	return true
}

//Method for room to add a new message
//func (r Room)NewMessage(u User, msgbody string)

//Room method to add a user to the room
func (r *Room) AddUser(u *User) {
	//r.Users = append(r.Users, u)
	u.JoinRoom(r)
}

//User method to add a room to the list of room that the user is part of
func (u *User) JoinRoom(r *Room) {
	//u.Rooms = append(u.Rooms, r)

}

// Removes the user from the room
func (r *Room) RemoveUser(u *User) {
	// for i, elem := range r.Users {
	// 	if reflect.DeepEqual(elem, u) {
	// 		r.Users = r.Users[:i+copy(r.Users[i:], r.Users[i+1:])]
	// 	}
	// }
	// r.NoPeople--

}

// Removes the room from the user
<<<<<<< HEAD
func (u *User)RemoveRoom(r *Room) {
	for i, elem := range u.Rooms {
		if reflect.DeepEqual(elem, r) {
			u.Rooms = u.Rooms[:i+copy(u.Rooms[i:], u.Rooms[i+1:])]
		}
	}
=======
func (u *User) RemoveRoom(r *Room) {
	// for i, elem := range u.Rooms {
	// 	if reflect.DeepEqual(elem, r) {
	// 		u.Rooms = u.Rooms[:i+copy(u.Rooms[i:], u.Rooms[i+1:])]
	// 	}
	// }
>>>>>>> 37c131a232a316ee80b37a72f3b95d0417e0ccec
}

func CreateUser(uname string, pass string, c net.Conn) *User {
	u := User{}
	u.UName = uname
	u.Pass = pass
	//u.Rooms = []Room{}
	//u.ActiveConn = c
	
	//Users = append(Users,&u)

	return &u
}

//Purpose: returns an array of the names of the rooms the user is in
//Use:    When the client software wants to list rooms, passing a name as an argument for joining a room, etc.
//Tested: NO
func (u *User) ShowRooms() []string {
	// var roomNames []string
	// for _, r := range u.Rooms {
	// 	roomNames = append(roomNames, r.Name)
	// }
	// return roomNames
	return nil
}

//Purpose: Creating a new room
//Use: To create a new chat room
//Tested: No
func CreateRoom(name string) *Room {
	r := Room{}
	r.Name = name
	// r.Users = []*User{}
	// r.Messages = []Message{}

	//Rooms = append(Rooms, &r)

	return &r
}

func CreateMessage(Uname string, text string, id int) *Message{
	m := Message{}
	m.Time = time.Now()
	m.Uname = Uname
	m.Body = text
	m.ID = id

	return &m
}

//Purpose: returns an array of the names of the users in the room
//Use: when the client or server wishes to know what users are in the room
//Tested: NO
func ShowUsers(r Room) []string {
	var users []string
	// for _, u := range r.Users {
	// 	fmt.Printf("%v\n", u.Uname)
	// 	users = append(users, u.Uname)
	// }
	return users
}

func GetUser(uname string) *User{
	user, exists := Users[uname]

	if exists {
		return &user
	} else {
		return nil
	}
}

func GetRoom(id int) *Room{
	room, exists := Rooms[id]

	if exists {
		return &room
	} else {
		return nil
	}
}

func (u *User)PostMsg(r *Room,text string){
	len := len(r.Messages)
	//oldmsg := msgs[len(msgs)-1] //senaste meddelandet

	newmsg := CreateMessage(u.Uname, text, len)

	r.Messages = append(r.Messages, newmsg)
}

func GetMessage(msgId int, r *Room) *Message{
	msg := r.Messages
	return msg[msgId]
}

func GetMessages(msgId int, r *Room) []*Message{
	msg := r.Messages
	msgs := msg[msgId:]
	return msgs
}

func DestroyUser(id string){

}

func DestroyRoom(id string){

}

func getOlderMessages(name string, room string) {

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
