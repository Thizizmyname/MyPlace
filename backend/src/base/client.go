package main

import (
    "net"
    "fmt"
    "bufio"
	"os"
)


func main() {
	fmt.Printf("Connecting to server..\n")
	//Skapar två nya variabler, conn(ection) och err(or)
	// err är nil om allt gick bra, annars inte nil
	conn, err := net.Dial("tcp","127.0.0.1:1337") //localhost och port 1337
    if err != nil { //om något gick fel 
		//Handle errors
		fmt.Printf("Error, %v", err)
    } else { //annars fortsätt
		fmt.Printf("Connection established to the server..\n")

		// Loopar för att skriva till servern
		for{
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("text to send: " )
			msg, _ := reader.ReadString('\n')
			fmt.Fprintf(conn,msg)

			


			// Koden nedanför är för att läsa tillbaka från servern till klienten
			text, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print("Message from server: ", text)
		
		}
	}

}
