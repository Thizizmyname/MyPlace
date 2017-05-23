package main

import (
    "net"
    "fmt"
    "bufio"
	"os/signal"
	"os"
	"log"
    "time"
    //"strings"
)



func main() {
/*
	fmt.Println("Dialing up.. Please wait")
	for x := range []uint64{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24}{
		fmt.Printf("\r%v%v%v\t%v%v               ", strings.Repeat("#",x),strings.Repeat("-",24-x), "Loading", uint64(float64(x)/0.24), "%")
		time.Sleep(50*time.Millisecond)
	}
	fmt.Printf("\nFINISHED DIALINGUP......\n")
  */
	fmt.Println("Connecting to server..\n")
	//Skapar två nya variabler, conn(ection) och err(or)
	// err är nil om allt gick bra, annars inte nil
	tcpLAddr,_ := net.ResolveTCPAddr("tcp","127.0.0.1:1337")
	conn, err := net.DialTCP("tcp",nil,tcpLAddr) //localhost och port 1337
    if err != nil { //om något gick fel 
		//Handle errors
		log.Fatal(err)
    } else { //annars fortsätt
		fmt.Printf("Connection established to the server..\n")
    fmt.Printf("Commands:\n%v%v%v%v",
              "1 - sign up\n",
              "2 - login\n",
              "3 - getRooms\n",
              "4 - createRoom (uses current time as name, so you can spam to get many rooms)\n")


		go readMsg(conn)

		c := make (chan os.Signal,1)
		signal.Notify(c, os.Interrupt)

		go disconnect(conn,c)
		

		//	WriteMsg

		for{
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("\rtext to send: \n" )
			msg, err := reader.ReadString('\n')
			if err != nil{
				log.Fatal(err)
			}
			switch msg{
      //Signup
      case "1\n":
        fmt.Printf("00{\"RequestID\":1,\"uname\":\"simon\",\"Pass\":\"hej\"}\n")
        fmt.Fprintf(conn,"00{\"RequestID\":1,\"UName\":\"simon\",\"Pass\":\"hej\"}\n")
      //Login
      case "2\n":
        fmt.Printf("01{\"RequestID\":1,\"UName\":\"simon\",\"Pass\":\"hej\"}\n")
        fmt.Fprintf(conn,"01{\"RequestID\":1,\"UName\":\"simon\",\"Pass\":\"hej\"}\n")
      //Create Room
      case "3\n":
        fmt.Printf("02{\"RequestID\":1,\"UName\":\"simon\",\"Pass\":\"hej\"}\n")
        fmt.Fprintf(conn,"02{\"RequestID\":1,\"UName\":\"simon\",\"Pass\":\"hej\"}\n")
      //getRooms
      case "4\n":
        fmt.Printf("08{\"RequestID\":1,\"UName\":\"simon\",\"Pass\":\"hej\"}\n")
        fmt.Fprintf(conn,"08{\"RequestID\":1,\"RoomName\":\"rum%v\",\"UName\":\"simon\"}\n", time.Now())
      default:
        fmt.Fprintf(conn,msg)
      }
			//fmt.Fprintf(conn,msg)			
			
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

func disconnect(conn net.Conn, channel chan os.Signal){

	sign := <-channel
	fmt.Println("\nDisconnecting... ", sign)

	os.Exit(3)
	// Skicka ett meddelande "Exit" till servern som handskas med den och kallar på en funktion som kör TCPClose på rätt connection.
	

}
