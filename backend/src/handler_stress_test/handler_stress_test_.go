package main

import (
	"io/ioutil"
	"io"
	"log"
	"fmt"
	"math/rand"
	"myplaceutils"
	"data"
	"handler"
	"requests_responses"
	"time"
	"os"
	"os/signal"
	"syscall"
)

var noUsers int = 0
var noRooms int = 0
var msgLengthIndicator int = 20
var stopFlag = false

func main() {
	InitLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	c := make(chan os.Signal, 2)
	c2 := make(chan bool)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stopFlag = true
		time.Sleep(time.Second)
		err := data.StoreDBs(myplaceutils.Users, myplaceutils.Rooms)
		if err!=nil {
			fmt.Printf("Error store db, %v\n", err)
		} else {
			fmt.Println("DBs successfully stored to file")
		}
		c2<-true
	}()


	rand.Seed(time.Now().Unix())
	myplaceutils.InitDBs()
	initLoggers(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	handlerChan := make(chan myplaceutils.HandlerArgs)
	go handler.ResponseHandler(handlerChan)
	defer close(handlerChan)

	for i := 0; !stopFlag; i++ {
		request := generateRequest(noUsers, noRooms)
		responseChan := make(chan requests_responses.Response, 20)
		handlerArgs := myplaceutils.HandlerArgs{request, responseChan}
		handlerChan <- handlerArgs
		response := <-responseChan

		fmt.Printf("%v.\n%T: %v\n%T: %v\n\n",
			i, request, request, response, response)

		if r, ok := response.(requests_responses.SignUpResponse); ok {
			if r.Result {
				noUsers++
			} else {
				panic("SignUp failed")
			}
		} else if _, ok := response.(requests_responses.CreateRoomResponse); ok {
			noRooms++
		}
	}
	<-c2
	os.Exit(1)
}

func generateRequest(noUsers int, noRooms int) requests_responses.Request {
	reqIndex := rand.Intn(12)

	switch reqIndex {
	case 0:
		return requests_responses.SignUpRequest{
			0, getNewUName(), "pass"}
	case 1:
		return requests_responses.SignInRequest{
			0, getExistingUName(), "pass"}
	case 2:
		return requests_responses.GetRoomsRequest{
			0, getExistingUName()}
	case 3:
		return requests_responses.GetRoomUsersRequest{
			0, getExistingRoomID()}
	case 4:
		return requests_responses.GetOlderMsgsRequest{
			0, getExistingRoomID(), 0}
	case 5:
		return requests_responses.GetNewerMsgsRequest{
			0, getExistingRoomID(), 0}
	case 6:
		return requests_responses.JoinRoomRequest{
			0, getExistingRoomID(), getExistingUName()}
	case 7:
		return requests_responses.LeaveRoomRequest{
			0, getExistingRoomID(), getExistingUName()}
	case 8:
		return requests_responses.CreateRoomRequest{
			0, "room", getExistingUName()}
	case 9:
		return requests_responses.PostMsgRequest{
			0, getExistingUName(), getExistingRoomID(), getRandomText()}
	case 10:
		return requests_responses.MsgReadRequest{
			0, 0, getExistingRoomID(), getExistingUName()}
	case 11:
		return requests_responses.SignOutRequest{
			0, getExistingUName()}
	default:
		panic("Bad reqIndex")
	}
}

func getNewUName() string {
	return fmt.Sprintf("user%v", noUsers)
}

//note: if 0 users, returns user100 (an non-existing user)
func getExistingUName() string {
	if noUsers == 0 {
		return "user100"
	} else {
		i := rand.Intn(noUsers)
		return fmt.Sprintf("user%v", i)
	}
}

//note: if 0 rooms, returns 100 (a non-existing roomID)
func getExistingRoomID() int {
	if noRooms == 0 {
		return 100
	} else {
		return rand.Intn(noRooms)
	}
}

func getRandomText() string {
        str := ""
        for i := 0; i < rand.Intn(msgLengthIndicator); i++ {
                str = fmt.Sprintf("%s %s ", str, "hello")
        }
        return str
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
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	myplaceutils.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	myplaceutils.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
