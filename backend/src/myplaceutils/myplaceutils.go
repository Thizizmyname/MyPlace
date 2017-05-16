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
func (u *User) LeaveRoom(r *Room) {
//	u.Rooms.Remove(r.Name)

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

// Purpose: Get the latest message from a room
// Argument: A room, which purpose is to get the latest message
// Returns a pointer to the latest message and the ID of the latest. If There are no messages in the room, returns nil and -1.
// Tested: Yes
func GetLatestMsg(room *Room) (*Message,int){
	maxID := -1

	if len(room.Messages) == 0 {
		return nil,maxID
	}

	
	for id,_ := range room.Messages{
		if id > maxID {
			maxID = id
		}
	}

	latestMsg := room.Messages[maxID]
	return latestMsg,maxID

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

func RemoveUsersOutgoingChannels(uname string, uOutChan chan requests_responses.Response) {
	user := GetUser(uname)

	if user == nil { return }

	for eR := user.Rooms.Front(); eR != nil; eR = eR.Next() {
		roomID := eR.Value.(int)
		room := GetRoom(roomID)

		for eOC := room.OutgoingChannels.Front(); eOC != nil; eOC = eOC.Next() {
			outChan := eOC.Value.(chan requests_responses.Response)
			if outChan == uOutChan {
				room.OutgoingChannels.Remove(eOC)
			}
		}
	}
}

//Adds an outgoing channel to the rooms list. If channel is already
//in the list, does nothing.
func (r *Room) AddOutgoingChannel(c chan requests_responses.Response) {
	if !outChanInUse(c, r.OutgoingChannels) {
		r.OutgoingChannels.PushBack(c)
	}
}

func outChanInUse(c chan requests_responses.Response, outChans *list.List) bool {
	for e := outChans.Front(); e != nil; e = e.Next() {
		outChan := e.Value.(chan requests_responses.Response)
		if c == outChan { return true }
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

// Purpose: Creates RoomInfo
// Argument: A room, a message and a username
// Returns: RoomInfo about a room
// Tested: No
func CreateRoomInfo(room *Room, msg *Message, username string) requests_responses.RoomInfo{
	msgInfo := CreateMsgInfo(room,msg,username)
	_,latestReadMsgID := GetLatestMsg(room)

	latestReadMsgID = -1 //Notera: ska returnera senast lästa meddelandet som har läst av användaren. För detta ska funka måste user-strukturen uppdateras, Ta bort denna när det har gjorts.

	roomInfo := requests_responses.RoomInfo{room.ID,room.Name,msgInfo,latestReadMsgID} 
	return roomInfo
}


// Purpose: Creates MessageInfo
// Argument: A room, a message and a username
// Returns: Info about a the latest message
// Tested: No
func CreateMsgInfo(room *Room, msg *Message, username string) requests_responses.MsgInfo {
	latestMsg,latestMsgID := GetLatestMsg(room)

	if len(room.Messages) == 0 {
		return requests_responses.MsgInfo{-1, -1, "", 0, ""}
	}

	
	msgInfo := requests_responses.MsgInfo{
		latestMsgID,
		room.ID,
		username,
		(latestMsg.Time).Unix(),
		latestMsg.Body }

	return msgInfo

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
