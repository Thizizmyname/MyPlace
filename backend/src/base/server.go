package main

import (
    "net"
    "fmt"
	//"time"
	"bufio"
//	"log"
	"reflect"
)

var connections []net.Conn

/*

*/
func main() {
	fmt.Printf("Starting server...\n")
	//Skapar två variabler listener och err
	//Om allt går bra kommer err vara nil, vilket gör det enkelt att testa
	//så att det gick bra
	tcpLAddr,_ := net.ResolveTCPAddr("tcp","127.0.0.1:1337") // För att få en tcp-address
	listener, err := net.ListenTCP("tcp",tcpLAddr)
	
	if err != nil { //Här testas om err inte är nil (alltså, har det blivit fel)
		fmt.Printf("Error starting server: %v\n", err)
	}
	for { //Denna loopen körs för evigt
		//Här skapas två nya variabler, connection och errs, samma som innan
		// men connection är varje anslutning som inkommer till server
		// dvs en socket
		connection, errs := listener.Accept() 
		if errs != nil { //Här testas igen om det blev något fel
			fmt.Printf("Connection Accept error: %v\n",errs)
		} else { //om inga fel inträffade, kan vi gå vidare
			connections = append(connections, connection)
			fmt.Printf("Connection established: %v\n", connection)
			go handleConnection(connection)
			go readMsg(connection)
			
		}
	}
}

/*
    Denna funktionen körs som en goroutine parallellt med main.
*/
func handleConnection(conn net.Conn) {
	fmt.Printf("Handle connection whattup: %v\nOpen connections:\n", conn)
	for _, conn := range connections{
		fmt.Println(conn)
	}
}

func readMsg(conn net.Conn){
	
	for{
		msg,err := bufio.NewReader(conn).ReadString('\n')
		if err == nil{

			
			fmt.Print("Message recieved: ", msg)
			for _, cons := range connections{
				sendMsg(cons,msg)
			}
		}else{

			// Tar just nu för givet att så fort det blir ett error så är det pga att en klient har stängt av
			// Tar bort klienten ur connections.
			removeConnection(conn)
			break
		}
	}
}

func sendMsg(conn net.Conn, msg string){

	fmt.Fprintf(conn,msg)
}


// Tar bort den specifika TCPAddressen i den globala slicen connections.
func removeConnection(conn net.Conn){
	for i, tcpAddr := range connections{
		if reflect.DeepEqual(tcpAddr,conn){
			fmt.Println("User ", tcpAddr, " has left the server")
			
			// Har ej skrivit denna men det är ett sätt för att ta bort önskat element
			connections = connections[:i+copy(connections[i:], connections[i+1:])] 
		}

	}
	
	
}


