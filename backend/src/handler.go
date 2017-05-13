package main

import (
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
		roomID := e.Value.(int)
		room := myplaceutils.GetRoom(roomID)
		room.OutgoingChannels.PushBack(responseChan)
	}

	return requests_responses.SignInResponse{requestID, true, ""}
}

func getRooms(request requests_responses.GetRoomsRequest) requests_responses.Response {
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.GetRoomsIndex, "not implemented yet"}
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
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.JoinRoomIndex, "not implemented yet"}
}

func leaveRoom(request requests_responses.LeaveRoomRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.LeaveRoomIndex, "not implemented yet"}
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
	newRoom.OutgoingChannels.PushBack(responseChan)

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
	return requests_responses.ErrorResponse{request.RequestID, requests_responses.MsgReadIndex, "not implemented yet"}
}

func signOut(request requests_responses.SignOutRequest, responseChan chan requests_responses.Response) requests_responses.Response {
	myplaceutils.RemoveUsersOutgoingChannels(request.UName, responseChan)
	return requests_responses.SignOutResponse{request.RequestID}
}
