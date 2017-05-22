package handler

import (

	"myplaceutils"
	"requests_responses"
	"net"
	"bufio"
	"fmt"
)



func ClientResponseHandler(conn net.Conn, clientResponseChannel chan requests_responses.Response) {
  myplaceutils.Info.Printf("Split client communication to clientResponseHandler\nConn: %v\n", conn)
  for args := range clientResponseChannel {
    responseString, err := requests_responses.ToResponseString(args)
    myplaceutils.Info.Printf("Parsed responseString: %v\n",responseString)
    if err!=nil{
      fmt.Fprintf(conn,"%v\n",err)
    }
    fmt.Fprintf(conn,"%v\n",responseString)
  }
}

func ClientHandler(conn net.Conn, clientChannel chan requests_responses.Response) {
  myplaceutils.Info.Println("Connection sent to clientHandler in go routine")
  go ClientResponseHandler(conn, clientChannel)
  var requestParsed myplaceutils.HandlerArgs
  var parseError error
  for {
    request ,err := bufio.NewReader(conn).ReadString('\n')
    if err!=nil {
      myplaceutils.Error.Println("User disconnected from the server")
      //TODO SignOutRequest
      break
    }
    myplaceutils.Info.Printf("New request: %v",request)
    requestParsed.IncomingRequest, parseError = requests_responses.FromRequestString(request)
    requestParsed.ResponseChannel = clientChannel
    if parseError==nil {
      myplaceutils.ResponseChannel <- requestParsed
    } else {
      myplaceutils.Error.Printf("Bad request from client: %v\n", parseError)
      //Har vi en default Bad request?
    }
  }
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
  var userRoomArray []requests_responses.RoomInfo
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
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.GetRoomUsersIndex, "not implemented yet"}
}

func getOlderMsgs(request requests_responses.GetOlderMsgsRequest) requests_responses.Response {
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.GetOlderMsgsIndex, "not implemented yet"}
}

func getNewerMsgs(request requests_responses.GetNewerMsgsRequest) requests_responses.Response {
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.GetNewerMsgsIndex, "not implemented yet"}
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

  if room == nil {
    return requests_responses.ErrorResponse{
      requestID,
      requests_responses.JoinRoomIndex,
      "Bad roomID"}
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
