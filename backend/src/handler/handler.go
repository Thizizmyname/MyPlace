package handler

import (
	"myplaceutils"
	"requests_responses"
	"net"
	"bufio"
	"fmt"
	"strings"
)



func ClientResponseHandler(conn net.Conn, clientResponseChannel chan requests_responses.Response) {
	myplaceutils.Info.Printf("Split client communication to clientResponseHandler\nConn: %v\n", conn)
	for args := range clientResponseChannel {
		responseString, err := requests_responses.ToResponseString(args)
		myplaceutils.Info.Printf("Parsed responseString: %v\n",responseString)
		if err!=nil {
			panic("Error while encrypting response")
		}

		_, writeErr := fmt.Fprintf(conn,"%v\n", responseString)

		if writeErr != nil {
			conn.Close()
			return
		}
	}
}

func ClientHandler(conn net.Conn, clientChannel chan requests_responses.Response) {
	myplaceutils.Info.Println("Connection sent to clientHandler in go routine")
	go ClientResponseHandler(conn, clientChannel)
	readBuf := make([]byte, myplaceutils.ConnReadMaxLength + 1)
	signedInUser := ""

	for {
		n, err := conn.Read(readBuf)

		if n > myplaceutils.ConnReadMaxLength {
			myplaceutils.Error.Printf("Request from client too large!")
			clientChannel <- requests_responses.ErrorResponse{-1, -1, "request too large, handler overwhelmed"}
		}

		if err != nil || n == 0 {
			myplaceutils.Info.Println("User disconnected from the server")
			if signedInUser != "" {
				sendSignOutRequest(-1, signedInUser, clientChannel)
			}
			conn.Close()
			return
		}

		reader := strings.NewReader(string(readBuf[:n]))
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			requestJson := scanner.Text()

			myplaceutils.Info.Printf("New request: %v", requestJson)
			request, parseError := requests_responses.FromRequestString(requestJson)

			if parseError != nil {
				myplaceutils.Error.Printf("Bad request from client: %v\n", parseError)
				clientChannel <- requests_responses.ErrorResponse{-1, -1, "parse error of request"}
				continue
			}

			if signInReq, ok := request.(requests_responses.SignInRequest); ok {
				if signedInUser == "" {
					signedInUser = signInReq.UName  //new user signing in
				} else {
					//trying to sign in without signing out previous user.
					//shouldn't happen.., still:
					sendSignOutRequest(signInReq.RequestID, signInReq.UName, clientChannel)
				}
			} else if _, ok := request.(requests_responses.SignOutRequest); ok {
				signedInUser = ""
			}

			var requestParsed myplaceutils.HandlerArgs
			requestParsed.IncomingRequest = request
			requestParsed.ResponseChannel = clientChannel
			myplaceutils.RequestChannel <- requestParsed
		}
	}
}

func sendSignOutRequest(reqID int, uname string, c chan requests_responses.Response) {
	args := myplaceutils.HandlerArgs{requests_responses.SignOutRequest{reqID, uname}, c}
	myplaceutils.RequestChannel <- args
}

func ResponseHandler(incomingChannel chan myplaceutils.HandlerArgs) {
	myplaceutils.Info.Println("Reached responseHandler")
	for args := range incomingChannel {
		request := args.IncomingRequest
		responseChan := args.ResponseChannel
		myplaceutils.Info.Printf("Handling request: %v\n", request)
		var response requests_responses.Response

		switch request.(type) {
		case requests_responses.SignUpRequest:
			myplaceutils.Info.Println("Matched request to signUp")
			response = signUp(request.(requests_responses.SignUpRequest))
		case requests_responses.SignInRequest:
			myplaceutils.Info.Println("Matched request to signIn")
			response = signIn(request.(requests_responses.SignInRequest), responseChan)
		case requests_responses.GetRoomsRequest:
			myplaceutils.Info.Println("Matched request to getRooms")
			response = getRooms(request.(requests_responses.GetRoomsRequest))
		case requests_responses.GetRoomUsersRequest:
			myplaceutils.Info.Println("Matched request to getRoomUsers")
			response = getRoomUsers(request.(requests_responses.GetRoomUsersRequest))
		case requests_responses.GetOlderMsgsRequest:
			myplaceutils.Info.Println("Matched request to GetOlderMsgs")
			response = getOlderMsgs(request.(requests_responses.GetOlderMsgsRequest))
		case requests_responses.GetNewerMsgsRequest:
			myplaceutils.Info.Println("Matched request to GetNewerMsgs")
			response = getNewerMsgs(request.(requests_responses.GetNewerMsgsRequest))
		case requests_responses.JoinRoomRequest:
			myplaceutils.Info.Println("Matched request to JoinRoom")
			response = joinRoom(request.(requests_responses.JoinRoomRequest), responseChan)
		case requests_responses.LeaveRoomRequest:
			myplaceutils.Info.Println("Matched request to LeaveRoom")
			response = leaveRoom(request.(requests_responses.LeaveRoomRequest), responseChan)
		case requests_responses.CreateRoomRequest:
			myplaceutils.Info.Println("Matched request to CreateRoom")
			response = createRoom(request.(requests_responses.CreateRoomRequest), responseChan)
		case requests_responses.PostMsgRequest:
			myplaceutils.Info.Println("Matched request to PostMsg")
			response = postMsg(request.(requests_responses.PostMsgRequest), responseChan)
		case requests_responses.MsgReadRequest:
			myplaceutils.Info.Println("Matched request to SignOut")
			response = msgRead(request.(requests_responses.MsgReadRequest))
		case requests_responses.SignOutRequest:
			myplaceutils.Info.Println("Matched request to SignOut")
			response = signOut(request.(requests_responses.SignOutRequest), responseChan)
		default:
			myplaceutils.Info.Println("Bad request type")
			panic("Nonexistent request type")
		}

		responseChan <- response
	}
}

func signUp(request requests_responses.SignUpRequest) requests_responses.Response {
	requestID := request.RequestID
	uname := request.UName
	pass := request.Pass

	if myplaceutils.UserExists(uname) {
		return requests_responses.SignUpResponse{requestID, false, "uname"}
	}

	if len(pass) < 3 {
		return requests_responses.SignUpResponse{requestID, false, "pass"}
	}

	myplaceutils.AddNewUser(uname, pass)
	response := requests_responses.SignUpResponse{request.RequestID, true, ""}

	return response
}

func signIn(request requests_responses.SignInRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	requestID := request.RequestID
	user := myplaceutils.GetUser(request.UName)
	pass := request.Pass

	if user == nil {
		return requests_responses.SignInResponse{requestID, false, "uname"}
	} else if pass != user.Pass {
		return requests_responses.SignInResponse{requestID, false, "pass"}
	}

	for e := user.Rooms.Front(); e != nil; e = e.Next() {
		room_ := e.Value.(myplaceutils.RoomIDAndLatestReadMsgID)
		room := myplaceutils.GetRoom(room_.RoomID)
		room.AddOutgoingChannel(responseChan)
	}

	return requests_responses.SignInResponse{requestID, true, ""}
}

func getRooms(request requests_responses.GetRoomsRequest) requests_responses.Response {
	/*
    1) Hitta user
    2) loopa över listan, ta e.Value.(typen)
    3)
  */
	userRoomArray := []requests_responses.RoomInfo{}
	user := myplaceutils.GetUser(request.UName)
	if user==nil{
		return requests_responses.ErrorResponse{request.RequestID, requests_responses.GetRoomsIndex, "User not found"}//TODO ändra denna till lämplig rad
	}
	for e := user.Rooms.Front(); e != nil; e = e.Next() {
		room := e.Value.(myplaceutils.RoomIDAndLatestReadMsgID)
		userRoom := myplaceutils.GetRoom(room.RoomID)

		userRoomArray = append(userRoomArray, myplaceutils.CreateRoomInfo(userRoom,user))
	}
	return requests_responses.GetRoomsResponse{request.RequestID, userRoomArray}
}

func getRoomUsers(request requests_responses.GetRoomUsersRequest) requests_responses.Response {
	id := request.RequestID
	roomId := request.RoomID

	room := myplaceutils.GetRoom(roomId)
	users := myplaceutils.ShowUsers(room)

	response := requests_responses.GetRoomUsersResponse{id, roomId, users }

	return response
}

func getOlderMsgs(request requests_responses.GetOlderMsgsRequest) requests_responses.Response {
	// Denna returenerar även det senaste lästa meddelandet
	id := request.RequestID
	roomID := request.RoomID
	msgID := request.MsgID // Senaste meddelandet som lästs
	NoMsgs := 20  // Anger hur många meddelande som ska hämtas
	room := myplaceutils.GetRoom(roomID)
	messages := []requests_responses.MsgInfo{}

	if room == nil {
		return requests_responses.ErrorResponse{
			id,
			requests_responses.GetOlderMsgsIndex,
			"bad roomID"}
	}

	if msgID > len(room.Messages) {
		return requests_responses.ErrorResponse{id,4,"MessageID does not exist"}
	}

	if NoMsgs > msgID {
		NoMsgs = msgID
	}

	for x := (msgID - NoMsgs); x < msgID; x++ {
		msg := room.Messages[x] // Meddelande x i rummet
		msginfo := myplaceutils.CreateMsgInfo(msg, roomID) // Skapar ett MsgInfo
		messages = append(messages,msginfo) // Lägger till MsgInfo till messages
	}
	response := requests_responses.GetOlderMsgsResponse{id,messages}
	return response
}

func getNewerMsgs(request requests_responses.GetNewerMsgsRequest) requests_responses.Response {

	requestID := request.RequestID
	roomID := request.RoomID
	msgID := request.MsgID
	getNoMsg := 10
	msgInfos := []requests_responses.MsgInfo{}

	room := myplaceutils.GetRoom(roomID)
	// Checks if the user has a msgID which is greater than the latestMsgID in the room

	if room == nil {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.GetNewerMsgsIndex,
			"bad roomID"}
	}

	if msgID > len(room.Messages) {
		return requests_responses.ErrorResponse{requestID,requests_responses.GetNewerMsgsIndex,"MessageID does not exist"}
	}

	if (getNoMsg + msgID) > len(room.Messages) {
		getNoMsg = len(room.Messages)
	}else{
		getNoMsg = msgID+getNoMsg +1 // inkluderar det sista meddelandet i rangen av taket från msgId
	}

	for x := (msgID+1) ; x < getNoMsg; x++ {
		msg := room.Messages[x]
		msgInfo := myplaceutils.CreateMsgInfo(msg,room.ID)
		msgInfos = append(msgInfos,msgInfo)
	}


	return requests_responses.GetNewerMsgsResponse{requestID,msgInfos}

}

func joinRoom(request requests_responses.JoinRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	// Vill uppdatera ett rum så att en user är medlem i det

	requestID := request.RequestID
	roomID := request.RoomID
	username := request.UName

	room := myplaceutils.GetRoom(roomID)
	user := myplaceutils.GetUser(username)

	if user == nil{
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.JoinRoomIndex,
			"There is no such user"}
	}

	if room == nil { // Kan inte skapa en respons utan att skapa en ha ett rum
		roomInfo := requests_responses.RoomInfo{}

		return requests_responses.JoinRoomResponse{
			requestID,
			roomInfo,
			false}
	}

	if myplaceutils.UserIsInRoom(username,room) {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.JoinRoomIndex,
			"User is already a member of the room"}
	}

	user.JoinRoom(room)
	room.AddOutgoingChannel(responseChan)

	roomInfo := myplaceutils.CreateRoomInfo(room,user)
	response := requests_responses.JoinRoomResponse{request.RequestID,roomInfo,true}

	return response
}

// Purpose: The user leaves a room
// Argument: a request, response channel
// Returns: a response
// Tested: Yes
func leaveRoom(request requests_responses.LeaveRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	// Vill uppdatera ett rum så att en user har lämnat i det
	requestID := request.RequestID
	roomID := request.RoomID
	username := request.UName

	room := myplaceutils.GetRoom(roomID)
	user := myplaceutils.GetUser(username)

	if user == nil{
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.LeaveRoomIndex,
			"There is no such user"}
	}

	if room == nil {

		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.LeaveRoomIndex,
			"Bad roomID"}
	}

	if !myplaceutils.UserIsInRoom(username,room) {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.LeaveRoomIndex,
			"There is no such user in the room"}
	}

	myplaceutils.RemoveUsersOutgoingChannels(user.UName,responseChan)
	user.LeaveRoom(room)
	return requests_responses.LeaveRoomResponse{requestID}
}

func createRoom(request requests_responses.CreateRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	requestID := request.RequestID
	roomName := request.RoomName
	user := myplaceutils.GetUser(request.UName)

	if user == nil {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.CreateRoomIndex,
			"no such user"}
	}

	newRoom := myplaceutils.AddNewRoom(roomName)
	user.JoinRoom(newRoom)
	newRoom.AddOutgoingChannel(responseChan)

	response := requests_responses.CreateRoomResponse{requestID, newRoom.ID, newRoom.Name}

	return response
}

func postMsg(request requests_responses.PostMsgRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	requestID := request.RequestID
	uname := request.UName
	roomID := request.RoomID
	body := request.Body
	room := myplaceutils.GetRoom(roomID)

	if room == nil {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.PostMsgIndex,
			"bad roomID"}
	} else if myplaceutils.UserIsInRoom(uname, room) == false {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.PostMsgIndex,
			"user not in room"}
	} else if len(body) == 0 || len(body) > myplaceutils.MsgMaxLength {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.PostMsgIndex,
			"bad msg length"}
	}

	msg := myplaceutils.AddNewMessage(uname, room, body)
	msgResp := requests_responses.MsgInfo{msg.ID, roomID, msg.UName, msg.Time.Unix(), msg.Body}

	requestIDToAllButSender := -1
	response := requests_responses.PostMsgResponse{requestIDToAllButSender, msgResp}

	for e := room.OutgoingChannels.Front(); e != nil; e = e.Next() {
		roomChan := e.Value.(chan requests_responses.Response)

		if roomChan != responseChan {
			roomChan <- response
		}
	}

	response.RequestID = requestID

	return response
}

func msgRead(request requests_responses.MsgReadRequest) requests_responses.Response {
	requestID := request.RequestID
	msgID := request.MsgID
	room := myplaceutils.GetRoom(request.RoomID)
	user := myplaceutils.GetUser(request.UName)

	if room == nil {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.MsgReadIndex,
			"bad roomID"}
	}
	if user == nil {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.MsgReadIndex,
			"bad userID"}
	}

	ok := user.SetLatestReadMsg(room, msgID)

	if ok == false {
		return requests_responses.ErrorResponse{
			requestID,
			requests_responses.MsgReadIndex,
			"bad msgID or user not in room"}
	}

	return requests_responses.MsgReadResponse{requestID}
}

func signOut(request requests_responses.SignOutRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	myplaceutils.RemoveUsersOutgoingChannels(request.UName, responseChan)
	return requests_responses.SignOutResponse{request.RequestID}
}
