package main

import "net"
import "bufio"
import "strings"
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
	noClientsToStartWith = 100
	maxSpawnDelay = 100
	maxRequestDelay = 1000
	noUsersPerRoom = 5
	noJoinedRoomsPerUser = 7
	printDelay = 2000
	msgLengthIndicator = 20
	maxNoSentChars = 4096
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
var out *log.Logger = log.New(os.Stdout, "", 0)

func main() {
	if len(os.Args) > 1 {
		dbg = log.New(os.Stdout, "\nEVENT:\n", 0)
	} else {
		dbg = log.New(ioutil.Discard, "\nEVENT:\n", 0)
	}

	rand.Seed(time.Now().Unix())
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
	conn, err := net.Dial("tcp", "127.0.0.1:1337")
	if err != nil {
		return false
	}

	createRoom := ""
	if index % noUsersPerRoom == 0 {
		createRoom = fmt.Sprintf("room%v", noRooms)
		noRooms++
	}

	clientWriterArgsChan := make(chan clientWriterArgs, 1)
	go clientSender(conn, clientWriterArgsChan, index, createRoom)
	go clientReceiver(conn, clientWriterArgsChan)
	dbg.Printf("New spawn, index:%v, create:%v", index, createRoom)

	return true
}


func clientSender(conn net.Conn, requestChan chan clientWriterArgs, index int, createRoom string) {
	uname := fmt.Sprintf("user%v", index)
	pass := "pass"
	joinRooms := getRandomRoomIDs(noRooms)

	requestID := 0
	signUpReq := requests_responses.SignUpRequest{
		requestID, uname, pass}
	sendAndSleep(conn, requestChan, signUpReq, requestID)

	requestID++
	signInReq := requests_responses.SignInRequest{
		requestID, uname, pass}
	sendAndSleep(conn, requestChan, signInReq, requestID)

	requestID++
	if createRoom != "" {
		createRoomReq := requests_responses.CreateRoomRequest{
			requestID, createRoom, uname}
		sendAndSleep(conn, requestChan, createRoomReq, requestID)
	}

	for _, roomID := range joinRooms {
		requestID++
		joinRoomReq := requests_responses.JoinRoomRequest{
			requestID, roomID, uname}
		sendAndSleep(conn, requestChan, joinRoomReq, requestID)
	}

	for ; ; {
		requestID++
		postRoom := joinRooms[rand.Intn(len(joinRooms))]
		body := getRandomMsgBody()
		postMsgReq := requests_responses.PostMsgRequest{
			requestID, uname, postRoom, body}
		sendAndSleep(conn, requestChan, postMsgReq, requestID)
	}
}

func sendAndSleep(conn net.Conn, requestChan chan clientWriterArgs, request requests_responses.Request, requestID int) {
	reqJ, _ := requests_responses.ToRequestString(request)
	dbg.Printf("Request sent: %v", reqJ)

	clientWriterArgs := clientWriterArgs{requestID, time.Now()}
	requestChan <- clientWriterArgs

	fmt.Fprintf(conn, "%s\n", reqJ)  //send request to server

	time.Sleep(time.Duration(rand.Intn(maxRequestDelay)) * time.Millisecond)
}

func clientReceiver(conn net.Conn, requestChan chan clientWriterArgs) {
	buf := make([]byte, maxNoSentChars + 1)
	requests := []clientWriterArgs{}

	for {
		n, err := conn.Read(buf)
		newReqs := newRequests(requestChan)
		requests = append(requests, newReqs...)

		if n > maxNoSentChars {
			panic("Too big response!")
		} else if n == 0 || err != nil {
			panic("Reading response failed")
		}

		reader := strings.NewReader(string(buf[:n]))
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			respJ := scanner.Text()
			dbg.Printf("Response recieved: %v", respJ)

			responseID := extractResponseID(respJ)
			if responseID == -1 { continue }

			responseTime := time.Now()
			var matchingReq clientWriterArgs
			requests, matchingReq = removeRequest(requests, responseID)
			requestTime := matchingReq.requestTime

			waitTime := responseTime.Sub(requestTime)
			waitTimeChan <- waitTime
		}

		if err := scanner.Err(); err != nil {
			panic("format error of incoming responses")
		}
	}
}

func gatherAndPrintStats() {
	maxWaitTime := 0.0
	noSamples := 0
	lastPrintTime := time.Now()

	noSamples_local := 0
	totalWaitTime_local := 0.0
	maxWaitTime_local := 0.0
	minWaitTime_local := 10000000.0
	out.Printf("Local data is gathered from the past %v milliseconds\n", printDelay)
	out.Println("noClients noSamples_global maxWaitTime_global noSamples_local maxWaitTime_local minWaitTime_local avgWaitTime_local")

	for waitTime_d := range waitTimeChan {
		waitTime := waitTime_d.Seconds()
		noSamples++

		if waitTime > maxWaitTime {
			maxWaitTime = waitTime
		}

		noSamples_local += 1
		totalWaitTime_local += waitTime

		if waitTime > maxWaitTime_local {
			maxWaitTime_local = waitTime
		}

		if waitTime < minWaitTime_local {
			minWaitTime_local = waitTime
		}

		t := time.Now()
		if lastPrintTime.Add(time.Duration(printDelay * time.Millisecond)).Before(t) {
			out.Printf("%v %v %.2f %v %.2f %.2f %.2f",
				noClients, noSamples, maxWaitTime, noSamples_local, maxWaitTime_local,
				minWaitTime_local, totalWaitTime_local / float64(noSamples_local))
			lastPrintTime = t

			noSamples_local = 0
			totalWaitTime_local = 0.0
			maxWaitTime_local = 0.0
			minWaitTime_local = 100000000.0
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
