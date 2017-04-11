package main

import (
    "net"
    "fmt"
	//"time"
	"bufio"
	"os"
)


/*

*/
func main() {
	fmt.Printf("Starting server...\n")
	//Skapar två variabler listener och err
	//Om allt går bra kommer err vara nil, vilket gör det enkelt att testa
	//så att det gick bra
	listener, err := net.Listen("tcp",":1337")
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
			fmt.Printf("Connection established: %v\n", connection)
			go handleConnection(connection)
			go readMsg(connection)
			go sendMsg(connection)
					
		}
	}
}

/*
    Denna funktionen körs som en goroutine parallellt med main.
*/
func handleConnection(conn net.Conn) {
	fmt.Printf("Handle connection whattup: %v\n", conn)
}

func readMsg(conn net.Conn){
	
	for{
	msg,_ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(" Message recivied: ", msg)
	}
}

func sendMsg(conn net.Conn){

	for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("text to send: ")
		msg, _ :=reader.ReadString('\n')
		fmt.Fprintf(conn,msg)
	}
}
