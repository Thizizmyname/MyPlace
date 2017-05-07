package myplaceutils

import (
    "fmt"
    "net"
    "time"
	"reflect"
	"log"
	"sync"
	"strings"
)


//Dessa loggers är till för att kunna anropas från alla programfiler.
var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
    connections []net.Conn
	GlobalUsers []*User
	GlobalRooms []*Room
)

type User struct {
	Uname string
	Pass string
	Rooms []*Room
	ActiveConn net.Conn
	mux sync.Mutex
  //token   string
}

type Room struct {
	Name     string
	NoPeople int
	Users    []*User
	Messages []Message // Den kanske ska innehålla pekare till meddelanden?
	mux sync.Mutex
	
}

type Message struct {
	Time  time.Time
	Uname string
	Body  string
	ID    string
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



//MÅSTE HA MUTEX LOCK I DENNA FUNKTIONEN
//MÅSTE KOLLA SÅ INTE newConnection REDAN FINNS I connections
func AddConnection(newConnection net.Conn) {
  connections = append(connections, newConnection)
}

//Kolla genom arrayen om den finns innan den försöker ta bort den
func RemoveConnection(connection net.Conn) bool{
  return true
}

//User method for binding the current connection to the user
func (u *User) BindConnection(c net.Conn) bool {
	u.ActiveConn = c
	return true
}

//Method for room to add a new message
//func (r Room)NewMessage(u User, msgbody string)

// Purpose: Room method to add a user to the room (Field Users)
// Use: When a user joins the room, this method updates so the user is a member of the room.
// Argument: The User that joins the room.
// Tested: Yes 
func (r *Room) AddUser(u *User) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Users = append(r.Users, u)
	r.NoPeople++

}

// Purpose: User method to add a room to the list of rooms that the user is part of (Field Rooms)
// Use: When the user joins a room this method updates the user. Argument: The room that the user want to join
// Argument: The room the user joins.
// Tested: Yes
func (u *User) JoinRoom(r *Room) {
	u.mux.Lock()
	defer u.mux.Unlock()
	u.Rooms = append(u.Rooms, r)

}

// Purpose:  Removes a user from the room (Field Users)
// Use: When a user wants to leave a room this method updates the specific room so the user isn't a member of the room.
// Argument: The user who leaves the room
// Return: True if it succeeds to remove a user from the room, else false
// Tested: Yes
func (r *Room) RemoveUser(u *User) bool {
	r.mux.Lock()
	defer r.mux.Unlock() //NOTIS: Varför behöver man inte ha lås här?
	for i, elem := range r.Users {
		if reflect.DeepEqual(elem, u) {
			r.Users = r.Users[:i+copy(r.Users[i:], r.Users[i+1:])]
			r.NoPeople--
			return true
		}
	}

	return false
}

// Purpose:  Removes the room from the user (Field Rooms)
// Use: When the user wants to leave a room this method updates the specific user
// Argument: The room the user leaves
// Return: True if it succeeds to remove a room from the user, else false
// Tested: Yes (a few tests)
func (u *User) RemoveRoom(r *Room) bool {
	r.mux.Lock()
	defer r.mux.Unlock() //NOTIS: Varför behöver man inte ha lås här?
	for i, elem := range u.Rooms {
		if reflect.DeepEqual(elem, r) {
			u.Rooms = u.Rooms[:i+copy(u.Rooms[i:], u.Rooms[i+1:])]
			return true
		}
	}

	return false
}

func CreateUser(uname string, pass string, c net.Conn) *User {
	u := User{}
	u.Uname = uname
	u.Pass = pass
	u.Rooms = []*Room{}
	u.ActiveConn = c

	GlobalUsers = append(GlobalUsers,&u)

	return &u
}

//Purpose: returns an array of the names of the rooms the user is in
//Use:    When the client software wants to list rooms, passing a name as an argument for joining a room, etc.
//Tested: NO
func (u *User) ShowRooms() []string {
	var roomNames []string
	for _, r := range u.Rooms {
		roomNames = append(roomNames, r.Name)
	}
	return roomNames
}

//Purpose: Creating a new room
//Use: To create a new chat room
//Tested: No
func CreateRoom(name string) *Room {
	r := Room{}
	r.Name = name
	r.NoPeople = 0
	r.Users = []*User{}
	r.Messages = []Message{}

	GlobalRooms = append(GlobalRooms, &r)
	
	return &r
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

func GetUser(id string) *User{
	
	for _,x := range GlobalUsers{
		if x.Uname == id{
			return x
		}
	}
	panic("can't find user")
}

func GetRoom(id string) *Room{

	for _,x := range GlobalRooms{
		if x.Name == id{
			return x
		}
	}
	panic("can't find room")	
}

func DestroyUser(id string){

}

func DestroyRoom(id string){
	
}

func getOlderMessages() {

}

func getNewerMessages() {

}

// Retunerar true om det namnet är ledigt
func CheckUsername(username string)bool {
	for _,elem := range GlobalUsers {
		if strings.Compare(elem.Uname,username) == 0 {
			return false
		}		
	}
	return true
}

// Check for illegal characters
func CheckCharacters(userName string) bool{
	return true
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
