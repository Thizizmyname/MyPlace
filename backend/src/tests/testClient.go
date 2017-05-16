package main

import (
    "net"
    "fmt"
    "bufio"
	"os/signal"
	"os"
	"log"
    "time"
    "strings"
    "strconv"
)



func main() {
	fmt.Println("Dialing up.. Please wait")
	for x := range []uint64{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24}{
		fmt.Printf("\r%v%v%v\t%v%v               ", strings.Repeat("#",x),strings.Repeat("-",24-x), "Loading", uint64(float64(x)/0.24), "%")
		time.Sleep(50*time.Millisecond)
	}
	fmt.Printf("\nFINISHED DIALINGUP......\n")
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


		go readMsg(conn)

		c := make (chan os.Signal,1)
		signal.Notify(c, os.Interrupt)

		go disconnect(conn,c)
		

		//	WriteMsg

		for{
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("\rCommand? (1-4): \n" )
			msg, err := reader.ReadString('\n')
			if err != nil{
				log.Fatal(err)
			}
			i, err := strconv.Atoi(msg[:len(msg)-1])
      if err!=nil {
        fmt.Printf("Error: %v\n",err)
      }
      switch i {
      case 1:
        signUp(conn)
      case 2:
        signIn(conn)
      case 3:
        getRooms(conn)
      case 4:
        getOlderMsgsRequest(conn)
      default:
        fmt.Println("Non valid entry.\nValid entries are:\n - 0 - signUp\n - 1 - signIn\n - 2 - getRooms\n - 3 - getRoomUsers\n - 4 - getOlderMsgs\n")
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

func signUp(conn net.Conn) {
  fmt.Fprintf(conn,"00{\"RequestID\":12345,\"UName\":\"user\",\"Pass\":\"pass\"}\n")
}


func signIn(conn net.Conn) {
  fmt.Fprintf(conn,"01{\"RequestID\":12345,\"UName\":\"user\",\"Pass\":\"pass\"}\n")
}


func getRooms(conn net.Conn) {
  fmt.Fprintf(conn,"02{\"RequestID\":12345,\"UName\":\"user\"}\n")
}


func getRoomUsers(conn net.Conn) {
  fmt.Println("LOW PRIO: Ej använd av frontend")
  fmt.Fprintf(conn,"03{\"RequestID\":12345,\"RoomID\":1}\n")
}

func getOlderMsgsRequest(conn net.Conn) {
  fmt.Fprintf(conn,"04{\"RequestID\":12345,\"RoomID\":1,\"MsgID\":2}\n")
}

