package main

import "time"
import "fmt"
import "requests_responses"
import "handler"
import "myplaceutils"
import "strconv"
import "math/rand"
import "log"
import "os"
import "io/ioutil"
import "io"

const (
	noClientsToStartWith = 1000
	maxSpawnDelay = 2
	maxRequestDelay = 5000
	noUsersPerRoom = 5
	noJoinedRoomsPerUser = 7
	printDelay = 2000
	msgLengthIndicator = 20
)

type clientWriterArgs struct {
	requestID int
	requestTime time.Time
}

var noClients int = 0
var noRooms int = 0
var waitTimeChan chan time.Duration = make(chan time.Duration, 20)
var handlerChan chan myplaceutils.HandlerArgs = make(chan myplaceutils.HandlerArgs)
var dbg *log.Logger
var out *log.Logger = log.New(os.Stdout, "\nOUTPUT: ", 0)

func main() {
	if len(os.Args) > 1 {
		dbg = log.New(os.Stdout, "\nEVENT:\n", 0)
	} else {
		dbg = log.New(ioutil.Discard, "\nEVENT:\n", 0)
	}

	myplaceutils.InitDBs()
	initLoggers(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	go handler.ResponseHandler(handlerChan)
	go gatherAndPrintStats()

	for i := 0; i < noClientsToStartWith; i++ {
		if startNewClient(noClients) {
			noClients++
		}
	}

	for ; ; {
		if startNewClient(noClients) {
			noClients++
		}

		time.Sleep(time.Duration(rand.Intn(maxSpawnDelay)) * time.Millisecond)
	}
}

func startNewClient(index int) bool {
	createRoom := ""
	if index % noUsersPerRoom == 0 {
		createRoom = fmt.Sprintf("room%v", noRooms)
		noRooms++
	}

	clientWriterArgsChan := make(chan clientWriterArgs, 1)
	responseChan := make(chan requests_responses.Response, 5)

	go clientSender(responseChan, clientWriterArgsChan, index, createRoom)
	go clientReceiver(responseChan, clientWriterArgsChan)
	dbg.Printf("New spawn, index:%v, create:%v", index, createRoom)

	return true
}


func clientSender(responseChan chan requests_responses.Response, requestChan chan clientWriterArgs, index int, createRoom string) {
	uname := fmt.Sprintf("user%v", index)
	pass := "pass"
	joinRooms := getRandomRoomIDs(noRooms)

	requestID := 0
	signUpReq := requests_responses.SignUpRequest{
		requestID, uname, pass}
	sendAndSleep(responseChan, requestChan, signUpReq, requestID)

	requestID++
	signInReq := requests_responses.SignInRequest{
		requestID, uname, pass}
	sendAndSleep(responseChan, requestChan, signInReq, requestID)

	requestID++
	if createRoom != "" {
		createRoomReq := requests_responses.CreateRoomRequest{
			requestID, createRoom, uname}
		sendAndSleep(responseChan, requestChan, createRoomReq, requestID)
	}

	for _, roomID := range joinRooms {
		requestID++
		joinRoomReq := requests_responses.JoinRoomRequest{
			requestID, roomID, uname}
		sendAndSleep(responseChan, requestChan, joinRoomReq, requestID)
	}

	for ; ; {
		requestID++
		postRoom := joinRooms[rand.Intn(len(joinRooms))]
		body := getRandomMsgBody()
		postMsgReq := requests_responses.PostMsgRequest{
			requestID, uname, postRoom, body}
		sendAndSleep(responseChan, requestChan, postMsgReq, requestID)
	}
}

func sendAndSleep(responseChan chan requests_responses.Response, requestChan chan clientWriterArgs, request requests_responses.Request, requestID int) {
	reqJ, _ := requests_responses.ToRequestString(request)
	dbg.Printf("Request sent: %v", reqJ)

	clientWriterArgs := clientWriterArgs{requestID, time.Now()}
	requestChan <- clientWriterArgs

	handlerArgs := myplaceutils.HandlerArgs{request, responseChan}
	handlerChan <- handlerArgs


	time.Sleep(time.Duration(rand.Intn(maxRequestDelay)) * time.Millisecond)
}

func clientReceiver(responseChan chan requests_responses.Response, requestChan chan clientWriterArgs) {
	requests := []clientWriterArgs{}

	for {
		response := <-responseChan
		newReqs := newRequests(requestChan)
		requests = append(requests, newReqs...)

		respJ, _ := requests_responses.ToResponseString(response)
		responseID := extractResponseID(respJ)
		if responseID == -1 { continue }

		dbg.Printf("Response recieved: %v", respJ)
		responseTime := time.Now()
		var matchingReq clientWriterArgs
		requests, matchingReq = removeRequest(requests, responseID)
		requestTime := matchingReq.requestTime

		waitTime := responseTime.Sub(requestTime)
		waitTimeChan <- waitTime
	}
}

func gatherAndPrintStats() {
	maxWaitTime := 0.0
	totalWaitTime := 0.0
	noSamples := 0
	lastPrintTime := time.Now()

	for waitTime_d := range waitTimeChan {
		waitTime := waitTime_d.Seconds()
		totalWaitTime += waitTime
		noSamples++

		if waitTime > maxWaitTime {
			maxWaitTime = waitTime
		}

		t := time.Now()
		if lastPrintTime.Add(time.Duration(printDelay * time.Millisecond)).Before(t) {
			dbg.Printf("New wait time: %0.2f. Total wait time: %0.2f", waitTime, totalWaitTime)
			out.Printf("no_clients:%v  no_responses:%v  max_wait:%.2f",
				noClients, noSamples, maxWaitTime)
			lastPrintTime = t
		}
	}
}

func getRandomRoomIDs(noRooms int) []int {
	ids := []int{}

	for i := 0; i < noJoinedRoomsPerUser; i++ {
		ids = append(ids, rand.Intn(noRooms))
	}

	return ids
}

func getRandomMsgBody() string {
	str := ""
	for i := 0; i < rand.Intn(msgLengthIndicator); i++ {
		str = fmt.Sprintf("%s %s ", str, "hello")
	}
	return str
}

func removeRequest(reqs []clientWriterArgs, respID int) ([]clientWriterArgs, clientWriterArgs) {
	for i, req := range reqs {
		if req.requestID == respID {
			reqs = append(reqs[:i], reqs[i+1:]...)
			return reqs, req
		}
	}

	panic("Response doesn't have a corresponding request")
}

func newRequests(c chan clientWriterArgs) []clientWriterArgs {
	args := []clientWriterArgs{}

	for {
		select {
		case x := <-c:
			if x.requestID != -1 {
				args = append(args, x)
			}
		default:
			return args
		}
	}
}

func extractResponseID(respJ string) int {
	firstColonIndex := -1
	for i, c := range respJ {
		if c == ':' {
			firstColonIndex = i
			break
		}
	}
	if firstColonIndex == -1 { panic("No colon in response") }

	firstCommaIndex := -1
	for i, c := range respJ {
		if c == ',' {
			firstCommaIndex = i
			break
		}
	}
	if firstColonIndex == -1 { panic("No comma in response") }

	id_str := respJ[firstColonIndex + 1 : firstCommaIndex]
	//fmt.Println(id_str)
	id, err := strconv.Atoi(id_str)
	if err != nil { panic(err) }

	return id
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
