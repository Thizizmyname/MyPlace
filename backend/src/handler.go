package main

import (
  "container/list"
  "myplaceutils"
  "requests_responses"
)


func handler(incomingChannel chan myplaceutils.HandlerArgs) {
	for args := range incomingChannel {
		request := args.IncomingRequest
		responseChan := args.ResponseChannel

		var response requests_responses.Response

		switch request.(type) {
		case requests_responses.SignUpRequest:
			response = signUp(request.(requests_responses.SignUpRequest))
		case requests_responses.SignInRequest:
			response = signIn(request.(requests_responses.SignInRequest), responseChan)
		case requests_responses.GetRoomsRequest:
			response = getRooms(request.(requests_responses.GetRoomsRequest))
		case requests_responses.GetRoomUsersRequest:
			response = getRoomUsers(request.(requests_responses.GetRoomUsersRequest))
		case requests_responses.GetOlderMsgsRequest:
			response = getOlderMsgs(request.(requests_responses.GetOlderMsgsRequest))
		case requests_responses.GetNewerMsgsRequest:
			response = getNewerMsgs(request.(requests_responses.GetNewerMsgsRequest))
		case requests_responses.JoinRoomRequest:
			response = joinRoom(request.(requests_responses.JoinRoomRequest), responseChan)
		case requests_responses.LeaveRoomRequest:
			response = leaveRoom(request.(requests_responses.LeaveRoomRequest), responseChan)
		case requests_responses.CreateRoomRequest:
			response = createRoom(request.(requests_responses.CreateRoomRequest), responseChan)
		case requests_responses.PostMsgRequest:
			response = postMsg(request.(requests_responses.PostMsgRequest), responseChan)
		case requests_responses.SignOutRequest:
			response = signOut(request.(requests_responses.SignOutRequest), responseChan)
		default:
			panic("Nonexistent request type")
		}

		responseChan <- response
	}
}

func signUp(request requests_responses.SignUpRequest) requests_responses.Response {
  //Example:
  //create a new user and update the db
  user := myplaceutils.User{
    request.UName,
    request.Pass,
    list.New()}

  //update the db:
  myplaceutils.Users[user.UName] = &user

  //create and return a response to the request
  response := requests_responses.SignUpResponse{request.RequestID, true, ""}

  return response

  //note: The request needs to be checked.. if UName is in use,
  //if pass ok etc. If error, the last string in the
  //SignUpResponse is set to error cause "user" or "pass" and
  //the bool is set to false, and the db isn't updated.
  //This stuff is defined in the request-response-interface.
}

func signIn(request requests_responses.SignInRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	return nil
}

func getRooms(request requests_responses.GetRoomsRequest) requests_responses.Response {
	//id := request.RequestID
	//name := request.UName

	//usr := myplaceutils.GetUser(name)
	//rooms := usr.Rooms

	
	
	return nil
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
  return nil
}

func getNewerMsgs(request requests_responses.GetNewerMsgsRequest) requests_responses.Response {
  return nil
}

func joinRoom(request requests_responses.JoinRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	return nil
}

func leaveRoom(request requests_responses.LeaveRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	return nil
}

func createRoom(request requests_responses.CreateRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	return nil
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

	for e := room.OutgoingChannels.Front(); e != nil; e.Next() {
		roomChan := e.Value.(chan requests_responses.Response)
		if roomChan != responseChan {
			roomChan <- response
		}
	}

	response.RequestID = requestID

	return response
}

func msgRead(request requests_responses.MsgReadRequest) requests_responses.Response {
  return nil
}

func signOut(request requests_responses.SignOutRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	return nil
}
