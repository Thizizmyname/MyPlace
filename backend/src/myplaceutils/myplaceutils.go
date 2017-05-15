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

var (
  //Dessa loggers är till för att kunna anropas från alla programfiler.
  Trace   *log.Logger
  Info    *log.Logger
  Warning *log.Logger
  Error   *log.Logger
  //connections []net.Conn
  Users UserDB
  Rooms RoomDB
  ResponseChannel chan HandlerArgs
)

const (
	MsgMaxLength = 1256
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
  Messages map[int]*Message //ID is key
  OutgoingChannels *list.List //[]chan requests_responses.Response
}

type Message struct {
	ID    int
	Time  time.Time
	UName string
	Body  string
}

type HandlerArgs struct {
  IncomingRequest requests_responses.Request
  ResponseChannel chan requests_responses.Response
}

type UserDB map[string]*User //UName is key
type RoomDB map[int]*Room //ID is key

func InitDBs() {
  Users = make(UserDB)
  Rooms = make(RoomDB)
}

/*
Funkar inte för att data.go är trasig
// Loads DB to global variables
func Initialize(){
  Users,Rooms,err = LoadDbs()
  //ska hantera erro på nåt sätt
}

//Store global variables to DB
func Terminate(){
  err := StoreDBs(Users,Rooms)
  //ska hantera error på nåt sätt
}
*/



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
	//u.JoinRoom(r)
}

//Updates dbs accordingly.
func (user *User) JoinRoom(room *Room) {
	if UserIsInRoom(user.UName, room) { return }

	user.Rooms.PushBack(room.ID)
	room.Users.PushBack(user.UName)
}

// Removes the user from the room
func (r *Room) RemoveUser(u *User) {
  // for i, elem := range r.Users {
  //   if reflect.DeepEqual(elem, u) {
  //     r.Users = r.Users[:i+copy(r.Users[i:], r.Users[i+1:])]
  //   }
  // }
  // r.NoPeople--

}

// Removes the room from the user
func (u *User) RemoveRoom(r *Room) {
  // for i, elem := range u.Rooms {
  //   if reflect.DeepEqual(elem, r) {
  //     u.Rooms = u.Rooms[:i+copy(u.Rooms[i:], u.Rooms[i+1:])]
  //   }
  // }
}

func CreateUser(uname string, pass string) *User {
	u := User{uname, pass, list.New()}
	return &u
}

//purpose: returns an array of the names of the rooms the user is in
//Use:    When the client software wants to list rooms, passing a name as an argument for joining a room, etc.
//Tested: NO
func (u *User) ShowRooms() []string {
  // var roomNames []string
  // for _, r := range u.Rooms {
  //   roomNames = append(roomNames, r.Name)
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

//Purpose: Create

//Purpose: returns an array of the names of the users in the room
//Use: when the client or server wishes to know what users are in the room
//Tested: NO
func ShowUsers(r Room) []string {
  var users []string
  // for _, u := range r.Users {
  //   fmt.Printf("%v\n", u.Uname)
  //   users = append(users, u.Uname)
  // }
  return users
}

//Purpose: Create a new user and add it to db, and return it.
func AddNewUser(uname string, pass string) *User {
	if UserExists(uname) {
		return nil
	} else {
		newUser := User{uname, pass, list.New()}
		Users[uname] = &newUser
		return &newUser
	}
}

//Purpose: Create a new room and add it to db, and return it.
func AddNewRoom(name string) *Room {
	newRoomID := findFreeRoomID()
	newRoom := Room{newRoomID, name, list.New(), make(map[int]*Message), list.New()}
	Rooms[newRoomID] = &newRoom
	return &newRoom
}

//Purpose: Create a new msg and add it to db, and return it.
func AddNewMessage(uname string, room *Room, body string) *Message {
	if UserIsInRoom(uname, room) == false {
		return nil
	}

	newMsgID := findFreeMsgID(room)
	newMsg := Message{newMsgID, time.Now(), uname, body}
	room.Messages[newMsg.ID] = &newMsg
	return &newMsg
}

func GetUser(uname string) *User{
  user, exists := Users[uname]

  if exists {
    return user
  } else {
    return nil
  }
}

func UserExists(uname string) bool {
	return GetUser(uname) != nil
}

func GetRoom(id int) *Room {
	room, exists := Rooms[id]

  if exists {
    return room
  } else {
    return nil
  }
}

func RoomExists(roomID int) bool {
	return GetRoom(roomID) != nil
}

func UserIsInRoom(uname string, room *Room) bool {
	unameList := room.Users

	for e := unameList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == uname {
			return true
		}
	}

	return false
}

func findFreeRoomID() int {
	maxID := -1

	for id, _ := range Rooms {
		if id > maxID { maxID = id }
	}

	return maxID + 1
}

func findFreeMsgID(r *Room) int {
	maxID := -1

	for id, _ := range r.Messages {
		if id > maxID { maxID = id }
	}

	return maxID + 1
}

func DestroyUser(id string){

}

func DestroyRoom(id string){

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
