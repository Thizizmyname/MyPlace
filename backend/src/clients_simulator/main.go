package main

import "time"
import "net"
import "fmt"
import "bufio"
import "requests_responses"
import "strconv"
import "strings"
import "math/rand"

const (
	noClientsToStartWith = 100
	spawnCount = 10
	spawnFreq = 1000
	requestFreq = 5000
	noRoomsRatioInv = 5
	noRoomsPerUser = 7
	maxNoSentChars = 1000000
	printFreq = 3
	msgLengthIndicator = 20
)

type clientWriterArgs struct {
	requestID int
	requestTime time.Time
}

var noClients int = 0
var noRooms int = 0
var waitTimeChan chan time.Duration

func main() {
	for i := 0; i < noClientsToStartWith; i++ {
		if startNewClient(noClients) {
			noClients++
		}
	}

	for {
		for i := 0; i < spawnCount; i++ {
			if startNewClient(noClients) {
				noClients++
			}
		}

		time.Sleep(time.Duration(spawnFreq) * time.Millisecond)
	}
}

func startNewClient(index int) bool {
	conn, err := net.Dial("tcp", "127.0.0.1:1337")

	if err != nil {
		return false
	}

	createRoom := ""
	if index % noRoomsRatioInv == 0 {
		createRoom = fmt.Sprintf("room%v", noRooms)
		noRooms++
	}

	clientWriterArgsChan := make(chan clientWriterArgs)
	go clientSender(conn, clientWriterArgsChan, index, createRoom)
	go clientReceiver(conn, clientWriterArgsChan)

	return true
}

func getRandomRoomIDs(noRooms int) []int {
	ids := []int{}

	for i := 0; i < noRoomsPerUser; i++ {
		ids = append(ids, rand.Intn(noRooms))
	}

	return ids
}

func getRandomMsgBody() string {
	str := ""
	for i := 0; i < rand.Intn(msgLengthIndicator); i++ {
		str = fmt.Sprintf("%s %s", str, "hello")
	}
	return str
}

func sendAndSleep(conn net.Conn, requestChan chan clientWriterArgs, request requests_responses.Request, requestID int) {
	jstr, err := requests_responses.ToRequestString(request)
	if err != nil {
		panic("error when parsing request")
	}

	clientWriterArgs := clientWriterArgs{requestID, time.Now()}
	requestChan <- clientWriterArgs
	fmt.Fprintf(conn, jstr)  //send request to server

	time.Sleep(time.Duration(requestFreq * time.Millisecond))
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

	for {
		requestID++
		postRoom := joinRooms[rand.Intn(len(joinRooms))]
		body := getRandomMsgBody()
		postMsgReq := requests_responses.PostMsgRequest{
			requestID, uname, postRoom, body}
		sendAndSleep(conn, requestChan, postMsgReq, requestID)
	}
}

func clientReceiver(conn net.Conn, requestChan chan clientWriterArgs) {
	buf := make([]byte, maxNoSentChars + 1)
	requests := []clientWriterArgs{}

	for {
		n, err := conn.Read(buf)

		if n > maxNoSentChars {
			panic("Too big response!")
		} else if n == 0 || err != nil {
			panic("Reading response failed")
		}

		reader := strings.NewReader(string(buf))
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			respJ := scanner.Text()
			responseTime := time.Now()
			responseID := extractResponseID(respJ)
			newReqs := newRequests(requestChan)
			requests = append(requests, newReqs...)

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
	//fmt.Println("no samples - avg wait time - max wait time")

	var maxWaitTime float64

	for {
		totalWaitTime := 0.0
		noSamples := 0

		for waitTime_d := range waitTimeChan {
			waitTime := waitTime_d.Seconds()
			totalWaitTime += waitTime
			noSamples++

			if waitTime > maxWaitTime {
				maxWaitTime = totalWaitTime
			}

			fmt.Printf("no_clients:%v:   wait:%v   max:%v\n",
				noClients, totalWaitTime / float64(noSamples), maxWaitTime)

			time.Sleep(time.Duration(printFreq * time.Second))
		}
	}
}

func removeRequest(reqs []clientWriterArgs, respID int) ([]clientWriterArgs, clientWriterArgs) {
	for i, req := range reqs {
		if req.requestID == respID {
			reqs = append(reqs[:i], reqs[i:1]...)
			return reqs, req
		}
	}

	panic("Response doen't have a corresponding request")
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

	id_str := respJ[firstColonIndex + 1:firstCommaIndex - 1]
	fmt.Println(id_str)
	id, err := strconv.Atoi(id_str)
	if err != nil { panic("") }

	return id
}
