package main

import (
	"io/ioutil"
	"io"
	"log"
	"os"
	"fmt"
	"myplaceutils"
	"handler"
	"requests_responses"
)

var no_users int = 0
var no_rooms int = 0

func main() {
	myplaceutils.InitDBs()
	initLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	handlerChan := make(chan myplaceutils.HandlerArgs)
	go handler.ResponseHandler(handlerChan)
	defer close(handlerChan)
	responseChan := make(chan requests_responses.Response)


	for i := 0; ; i++ {
		request := generateRequest()
		handlerArgs := myplaceutils.HandlerArgs{request, responseChan}
		handlerChan <- handlerArgs
		response := <-responseChan

		fmt.Printf("%v.\nRequest:  %v\nResponse: %v\n\n",
			i, request, response)
	}
}

func generateRequest() {

}



func initLoggers(
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
