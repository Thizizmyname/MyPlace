package main

import (
    "net"
    "fmt"
  //"time"
    "bufio"
//  "log"
//  "reflect"
  //"data" //importera denna när vi ska implementera databasen
//  "myplaceutils"
)


/*
Denna funktion ska användas för att ta emot en response från en klient
Just nu ligger det ett anrop för att den ska parsa meddelandet som den får in, men det kan bytas ut det med.
Den är inte färdigställd alls så allt från argument till bodyn kan förändras
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
  return ""
}

//Denna funktion ska användas för att skicka en response till en klient
//Den är inte färdigställd alls så allt från argument till bodyn kan förändras
func respondClient(clientConnection net.Conn, msg string){
  fmt.Fprintf(clientConnection,msg)
}


func handler(connection net.Conn, gochan chan string){
  //Placeholder
  for {
    
  }
}


func signup(username string, password string) {

}


func signin(username string, password string) /*TOKEN*/ string{
  //throw error
  return ""
}

func getRooms(username string /*Kanske byta mot user-type*/) []string{
  return []string{"h","a","ha"}
}

func getOlderMessages(roomID string, msgID string) []string /*kanske någon annan returntype*/ {
  //
  return []string{"he","he","he"}
}

func getNewerMessages(roomID string, msgID string) []string{
  //
  return []string{"ho","ho","ho"}
}

func joinRoom(roomId string, username string) {
// Exempel på hur det kan se ut.	
//	workingUser = getUser(username)
//	workingRoom = getRoom(roomId)
//	workingUser.JoinRoom(workingRoom)
	
}

func leaveRoom(username string, roomId string) {
  //
}

func createRoom(username string, roomName string) {
  //
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



