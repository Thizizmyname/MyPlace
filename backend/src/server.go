package main

import (
    "net"
    "io"
    "io/ioutil"
    "log"
    "os"
    "myplaceutils"
    "requests_responses"
    "data"
    "handler"
)


var connections []net.Conn

func listenLoop(listener net.Listener) {
  for { //Denna loopen körs för evigt
    //Här skapas två nya variabler, connection och errs, samma som innan
    // men connection är varje anslutning som inkommer till server
    // dvs en socket
    //channel := make(chan string)
    newConnection, errs := listener.Accept()
    if errs != nil { //Här testas igen om det blev något fel
      myplaceutils.Error.Printf("Connection Accept error: %v\n",errs)
    } else { //om inga fel inträffade, kan vi gå vidare
      connections = append(connections, newConnection)
      myplaceutils.Info.Printf("Connection established: %v\n", newConnection)
      clientChannel := make(chan requests_responses.Response, 8)
      go handler.ClientHandler(newConnection, clientChannel)
    }
  }
}

//Initialize loggers
func InitLoggers(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer,
    ) {
    myplaceutils.Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    myplaceutils.Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    myplaceutils.Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    myplaceutils.Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

//Om servern trycker ctrl+c ska denna anropas.
func disconnectAll() {
  for _,conn := range connections {
    conn.Close()
  }
}

func main() {
  //Title
  myplaceutils.PrintTitle()
  //Initializing
  //Loggers
  log.Println("Initializing loggers")
  InitLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
  //Users and Rooms
  myplaceutils.Info.Println("Loading database..")
  var loadError error
  myplaceutils.Users, myplaceutils.Rooms, loadError = data.LoadDBs()
  if loadError!=nil {
    myplaceutils.Error.Printf("Error loading database, continuing with empty databases.\nError message: %v\n",loadError)
  }
  err := data.StoreDBs(myplaceutils.Users, myplaceutils.Rooms)
  if err!=nil {
    myplaceutils.Error.Println(err)
  }
  log.Println("Initialization complete\n-------------------------")
  //255 känns som en lämplig size så länge MAGIC NUMBER
  myplaceutils.ResponseChannel = make(chan myplaceutils.HandlerArgs, 255)
  //Initialize RequestChannel, TODO döp om skiten till RequestChannel från ResponseChannel

	go handler.ResponseHandler(myplaceutils.ResponseChannel)
  myplaceutils.Info.Println("Creating a listener")

  tcpAddress,_ := net.ResolveTCPAddr("tcp",":1337")
  listener, err := net.ListenTCP("tcp",tcpAddress)

  if err != nil {
    log.Fatalf("myplaceutils.Error starting server: %v\n", err)
  } else {
    listenLoop(listener)
  }

}
