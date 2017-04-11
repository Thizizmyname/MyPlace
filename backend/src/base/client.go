package main

import (
    "net"
    "fmt"
    "bufio"
	"os"
	"log"
)


func main() {
	fmt.Printf("Connecting to server..\n")
	//Skapar två nya variabler, conn(ection) och err(or)
	// err är nil om allt gick bra, annars inte nil
	tcpLAddr,_ := net.ResolveTCPAddr("tcp","127.0.0.1:1337")
	//tcpRAddr,_ := net.ResolveTCPAddr("tcp","130.238.246.194:1337")
	conn, err := net.DialTCP("tcp",nil,tcpLAddr) //localhost och port 1337
    if err != nil { //om något gick fel 
		//Handle errors
		log.Fatal(err)
    } else { //annars fortsätt
		fmt.Printf("Connection established to the server..\n")

		// Loopar för att skriva till servern
		for{
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("text to send: " )
			msg, err := reader.ReadString('\n')

			if err != nil{
				log.Fatal(err)
			}
			
			fmt.Fprintf(conn,msg)
			
			


			// Koden nedanför är för att läsa tillbaka från servern till klienten
			text, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil{
				
			}

			fmt.Print("Message from server: ", text)

			
			
		}
	}

}
