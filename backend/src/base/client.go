package main

import (
    "net"
    "fmt"
    "bufio"
	"os"
	"log"
    "time"
    "strings"
)



func main() {
  fmt.Println("Dialing up.. Please wait")
  for x := range []uint64{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24}{
    fmt.Printf("\r%v%v%v\t%v%v               ", strings.Repeat("#",x),strings.Repeat("-",24-x), "Loading", uint64(float64(x)/0.24), "%")
    time.Sleep(500*time.Millisecond)
  }
  fmt.Printf("\nFINISHED DIALINGUP......\n")
	fmt.Println("Connecting to server..\n")
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
    go readMsg(conn)

		for{
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("\rtext to send: \n" )
			msg, err := reader.ReadString('\n')
      if err != nil{
				log.Fatal(err)
			}
			
			fmt.Fprintf(conn,msg)
			
			
    }

	}

}
func sleeper() {
  for x := range []uint64{0,1,2,3,4,5,6,7,8,9}{
    fmt.Printf("\r%v", '#'*x)
    time.Sleep(1)
  }
}

func readMsg(conn net.Conn) {
  for { 
			// Koden nedanför är för att läsa tillbaka från servern till klienten
			text, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil{
				log.Fatal(err)
			}
      
			fmt.Printf("\rMessage from server: %v\n", text)

			
			
		}
}
