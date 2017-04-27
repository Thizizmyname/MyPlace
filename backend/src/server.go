package main

import (
    "net"
    "io"
    "io/ioutil"
    "log"
    "os"
    "myplaceutils"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

var connections []net.Conn

func listenLoop(listener net.Listener) {
  for { //Denna loopen körs för evigt
    //Här skapas två nya variabler, connection och errs, samma som innan
    // men connection är varje anslutning som inkommer till server
   // dvs en socket
    channel := make(chan string)
    newConnection, errs := listener.Accept()
    if errs != nil { //Här testas igen om det blev något fel
      Error.Printf("Connection Accept error: %v\n",errs)
    } else { //om inga fel inträffade, kan vi gå vidare
      connections = append(connections, newConnection)
      Info.Printf("Connection established: %v\n", newConnection)
      Info.Printf("New channel\n", channel)
      go handler(newConnection, channel)
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
    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
  //Title
  myplaceutils.PrintTitle()
  //Initialize loggers
  log.Println("Initializing loggers")
  InitLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
  log.Println("Initialization complete")

  Info.Println("Creating a listener")

  tcpLAddr,_ := net.ResolveTCPAddr("tcp","127.0.0.1:1337")
  listener, err := net.ListenTCP("tcp",tcpLAddr)

  if err != nil {
    log.Fatalf("Error starting server: %v\n", err)
  } else {
    listenLoop(listener)
  }
  
}
