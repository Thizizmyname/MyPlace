package main

import (
	"container/list"
	"myplaceutils"
	"requests_responses"
)


func handler(incomingChannel chan myplaceutils.HandlerArgs) {
	for args := range incomingChannel {
		var request requests_responses.Request = args.IncomingRequest
		var response requests_responses.Response

		switch request.(type) {
		case requests_responses.SignUpRequest:
			response = signUp(request.(requests_responses.SignUpRequest))
		case requests_responses.SignInRequest:
			response = signIn(request.(requests_responses.SignInRequest))
		case requests_responses.GetRoomsRequest:
			response = getRooms(request.(requests_responses.GetRoomsRequest))
		case requests_responses.GetRoomUsersRequest:
			response = getRoomUsers(request.(requests_responses.GetRoomUsersRequest))
		case requests_responses.GetOlderMsgsRequest:
			response = getOlderMsgs(request.(requests_responses.GetOlderMsgsRequest))
		case requests_responses.GetNewerMsgsRequest:
			response = getNewerMsgs(request.(requests_responses.GetNewerMsgsRequest))
		case requests_responses.JoinRoomRequest:
			response = joinRoom(request.(requests_responses.JoinRoomRequest))
		case requests_responses.LeaveRoomRequest:
			response = leaveRoom(request.(requests_responses.LeaveRoomRequest))
		case requests_responses.CreateRoomRequest:
			response = createRoom(request.(requests_responses.CreateRoomRequest))
		case requests_responses.PostMsgRequest:
			response = postMsg(request.(requests_responses.PostMsgRequest))
		case requests_responses.SignOutRequest:
			response = signOut(request.(requests_responses.SignOutRequest))
		default:
			panic("Nonexistent request type")
		}

		args.ResponseChannel <- response
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
	myplaceutils.Users[user.UName] = user

	//create and return a response to the request
	response := requests_responses.SignUpResponse{request.RequestID, true, ""}

	return response

	//note: The request needs to be checked.. if UName is in use,
	//if pass ok etc. If error, the last string in the
	//SignUpResponse is set to error cause "user" or "pass" and
	//the bool is set to false, and the db isn't updated.
	//This stuff is defined in the request-response-interface.
}

func signIn(request requests_responses.SignInRequest) requests_responses.Response {
	return nil
}

func getRooms(request requests_responses.GetRoomsRequest) requests_responses.Response {
	return nil
}

func getRoomUsers(request requests_responses.GetRoomUsersRequest) requests_responses.Response {
	return nil
}

func getOlderMsgs(request requests_responses.GetOlderMsgsRequest) requests_responses.Response {
	return nil
}

func getNewerMsgs(request requests_responses.GetNewerMsgsRequest) requests_responses.Response {
	return nil
}

func joinRoom(request requests_responses.JoinRoomRequest) requests_responses.Response {
	return nil
}

func leaveRoom(request requests_responses.LeaveRoomRequest) requests_responses.Response {
	return nil
}

func createRoom(request requests_responses.CreateRoomRequest) requests_responses.Response {
	return nil
}

func postMsg(request requests_responses.PostMsgRequest) requests_responses.Response {
	return nil
}

func msgRead(request requests_responses.MsgReadRequest) requests_responses.Response {
	return nil
}

func signOut(request requests_responses.SignOutRequest) requests_responses.Response {
	return nil
}
