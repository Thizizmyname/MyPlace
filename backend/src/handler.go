package main

import (
  "container/list"
  "myplaceutils"
  "requests_responses"
  "net"
  "bufio"
)


func clientHandler(conn net.Conn, clientChannel chan requests_responses.Response) {
  myplaceutils.Info.Println("Connection sent to clientHandler in go routine")
  var requestParsed myplaceutils.HandlerArgs
  var parseError error
  for {
    request ,err := bufio.NewReader(conn).ReadString('\n')
    if err!=nil {
      myplaceutils.Error.Println("User disconnected from the server")
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

func responseHandler(incomingChannel chan myplaceutils.HandlerArgs) {
  myplaceutils.Info.Println("Reached responseHandler")
  for args := range incomingChannel {
    var request requests_responses.Request = args.IncomingRequest
    var response requests_responses.Response

    switch request.(type) {
    case requests_responses.SignUpRequest:
      myplaceutils.Info.Println("Matched request to signUp")
      response = signUp(request.(requests_responses.SignUpRequest))
    case requests_responses.SignInRequest:
      myplaceutils.Info.Println("Matched request to signIn")
      response = signIn(request.(requests_responses.SignInRequest))
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
      response = joinRoom(request.(requests_responses.JoinRoomRequest))
    case requests_responses.LeaveRoomRequest:
      myplaceutils.Info.Println("Matched request to LeaveRoom")
      response = leaveRoom(request.(requests_responses.LeaveRoomRequest))
    case requests_responses.CreateRoomRequest:
      myplaceutils.Info.Println("Matched request to CreateRoom")
      response = createRoom(request.(requests_responses.CreateRoomRequest))
    case requests_responses.PostMsgRequest:
      myplaceutils.Info.Println("Matched request to PostMsg")
      response = postMsg(request.(requests_responses.PostMsgRequest))
    case requests_responses.SignOutRequest:
      myplaceutils.Info.Println("Matched request to SignOut")
      response = signOut(request.(requests_responses.SignOutRequest))
    default:
      myplaceutils.Info.Println("Bad request type")
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
