package main

import (
    "net"
    "fmt"
	//"time"
    "bufio"
	//  "log"
	//  "reflect"
	//"data" //importera denna när vi ska implementera databasen
	"myplaceutils"
)
/*
Denna funktion ska användas för att ta emot en response från en klient
Just nu ligger det ett anrop för att den ska parsa meddelandet som den får in, men det kan bytas ut det med.
Den är inte färdigställd alls så allt från argument till bodyn kan förändras
Mer rimligt att detta funkar som en void-funktion och skicka vidare stringen till parsern istället för att returnera en string
*/
func recieveRequest(clientConnection net.Conn) string{
	for {
		clientRequestRaw,err := bufio.NewReader(clientConnection).ReadString('\n')
		if err!=nil {
			//clientRequestParsed := PARSE_FUNCTION_FROM_MARTIN(clientRequestRaw)
			clientRequestParsed := clientRequestRaw //TA BORT DENNA RADEN NÄR PARSE_FUNCTION_FROM_MARTIN finns
			return clientRequestParsed
		} else {
			//Kanske lämpligt att skicka tillbaka ett svar att det gick dåligt
		}
	}
}

//Denna funktion ska användas för att skicka en response till en klient
//Den är inte färdigställd alls så allt från argument till bodyn kan förändras
func respondClient(clientConnection net.Conn, msg string){
	fmt.Fprintf(clientConnection,msg)
}


func handler(connection net.Conn, gochan chan string){
	myplaceutils.Info.Printf("Reached handler, connection %v\n",connection)
	//Placeholder
	for {
		request ,err := bufio.NewReader(connection).ReadString('\n')
		if err==nil {
			myplaceutils.Info.Printf("Request: %v", request)
			//parsedRequest := myplaceutils.Parse(request) //DESSA TVÅ RADER SKA BYTAS UT
			//handleRequest(parsedRequest) 
		} else {
			myplaceutils.Error.Println("User disconnected from the server")
			
			myplaceutils.RemoveConnection(connection)
			break
		}
	}
}

/*
00.signup
args: uname, pass
response: -
note: error if uname is taken/ pass to short/ illegal characters/ ...
side-effect: updates users_db
*/
func SignUp(username string, password string, conn net.Conn) (bool){
	if myplaceutils.CheckCharacters(username) {
		
		if myplaceutils.CheckUsername(username) {
		//	fmt.Printf("The user name isn't taken")
			myplaceutils.CreateUser(username,password,conn)

		}else{
		//	fmt.Printf("The Username is allready taken. Chose a new username\n")
			return false
		}
		
	}else{
	//	fmt.Printf("Illegal characters")
		return false

	}
	return true
}

/*
01.signin
args: uname, pass
response: -
note: error if uname not in use/ incorrect pass
side-effect: new messages from all rooms that the client has joined start being pushed to client
*/
func signin(username string, password string) /*TOKEN*/ string{
	//throw error
	return ""
}

func getRooms(username string /*Kanske byta mot user-type*/) []string{
	
	return []string{"h","a","ha"}
}

func getOlderMessages(roomID string, msgID string) []string /*kanske någon annan returntype*/ {
	// ska användarna eller rummen ha channel i sig?
	//ME
	return []string{"he","he","he"}
}

func getNewerMessages(roomID string, msgID string) []string{

	
	//
	//ME
	// vet inte riktigt hur denna ska funka.
	return []string{"ho","ho","ho"}
}

func joinRoom(roomId string, username string){
	workingUser := myplaceutils.GetUser(username)
	workingRoom := myplaceutils.GetRoom(roomId)
	workingUser.JoinRoom(workingRoom) // Hur vet den vilken av joinRoom och JoinRoom(myplaceutils)?
	
}

func leaveRoom(username string, roomId string) {
	workingUser := myplaceutils.GetUser(username)
	workingRoom := myplaceutils.GetRoom(roomId)
	workingUser.RemoveRoom(workingRoom)

}

func createRoom(username string, roomName string) *myplaceutils.Room {	
	r := myplaceutils.CreateRoom(roomName)
	u := myplaceutils.GetUser(username)
	r.AddUser(u) 
	return r
}

func postMessage(username string, roomID string, text string){
	//
}

func messageRead(username string, roomID string, msgID string) {
	//
}



func deleteUser(username string) bool{
	return true
}



