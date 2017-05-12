package main

import (
	"testing"
	"os"
	"reflect"
	"myplaceutils"
	"requests_responses"
)

func TestMain(m *testing.M) {
	myplaceutils.InitDBs()
	retCode := m.Run()
	os.Exit(retCode)
}

func executeAndTestResponse(t *testing.T, request requests_responses.Request, expectedResponse requests_responses.Response) {
	handlerChan := make(chan myplaceutils.HandlerArgs)
	go handler(handlerChan) //now handler is waiting for requests
	defer close(handlerChan)

	responseChan := make(chan requests_responses.Response)
	handlerArgs := myplaceutils.HandlerArgs{request, responseChan}

	handlerChan <- handlerArgs //send args to handler
	response := <-responseChan

	if r, ok := response.(requests_responses.PostMsgResponse); ok {
		r2 := expectedResponse.(requests_responses.PostMsgResponse)
		if r.RequestID != r2.RequestID ||
			r.Msg.MsgID != r2.Msg.MsgID ||
			r.Msg.RoomID != r2.Msg.RoomID ||
			r.Msg.UName != r2.Msg.UName ||
			r.Msg.Body != r2.Msg.Body {

			t.Errorf("request: %v\nresponse: %v\nactual response: %v\nexpected response:%v",
			reflect.TypeOf(r), reflect.TypeOf(r2), r, r2)

		}
	} else if response != expectedResponse {
		t.Errorf("request: %v\nresponse: %v\nactual response: %v\nexpected response:%v",
			reflect.TypeOf(request),
			reflect.TypeOf(response),
			response,
			expectedResponse)
	}
}

// func executeAndTestResponsePostMsg(t *testing.T, request requests_responses.PostMsgRequest, r2 requests_responses.PostMsgResponse, room *myplaceutils.Room) {
// 	handlerChan := make(chan myplaceutils.HandlerArgs)
// 	go handler(handlerChan) //now handler is waiting for requests
// 	defer close(handlerChan)

// 	responseChan := make(chan requests_responses.Response)
// 	handlerArgs := myplaceutils.HandlerArgs{request, responseChan}

// 	handlerChan <- handlerArgs //send args to handler
// 	r := (<-responseChan).(requests_responses.PostMsgResponse)

// 	if r.RequestID != r2.RequestID ||
// 		r.Msg.MsgID != r2.Msg.MsgID ||
// 		r.Msg.RoomID != r2.Msg.RoomID ||
// 		r.Msg.UName != r2.Msg.UName ||
// 		r.Msg.Body != r2.Msg.Body {

// 		//t.Error(reflect.TypeOf(request))
// 	}

// 	for e := room.OutgoingChannels.Front(); e != nil; e = e.Next() {
// 		uChan := e.Value.(chan requests_responses.Response)
// 		r := (<-responseChan).(requests_responses.PostMsgResponse)
// 		if (uChan == responseChan) {
// 			if r.RequestID != r2.RequestID ||
// 				r.Msg.MsgID != r2.Msg.MsgID ||
// 				r.Msg.RoomID != r2.Msg.RoomID ||
// 				r.Msg.UName != r2.Msg.UName ||
// 				r.Msg.Body != r2.Msg.Body {

// 				t.Error(reflect.TypeOf(request))
// 			}
// 		} else {
// 			if r.RequestID != -1 ||
// 				r.Msg.MsgID != r2.Msg.MsgID ||
// 				r.Msg.RoomID != r2.Msg.RoomID ||
// 				r.Msg.UName != r2.Msg.UName ||
// 				r.Msg.Body != r2.Msg.Body {

// 				t.Error(reflect.TypeOf(request))
// 			}
// 		}
// 	}
// }


/*
00.signup
args: uname, pass
response: -
note: error if uname is taken/ pass to short/ illegal characters/ ...
side-effect: updates users_db
*/
func TestSignUp(t *testing.T) {
	//Example:
	//Create request and expected response to this request
	req := requests_responses.SignUpRequest{12345, "laban", "pass"}
	resp := requests_responses.SignUpResponse{12345, true, ""}

	//Reset the global dbs (rooms and users)
	myplaceutils.InitDBs()

	//The following call sends request (req) to the handler
	//and tests if the expected response (resp) is returned.
	//The dbs are updated by the handler according to the request.
	executeAndTestResponse(t, req, resp)

	//If expected response wasn't returned by handler, an error is
	//already raised.
	//So finally, check if the dbs were updated as expected or else
	//raise an error.
	//Does a "laban" exists in userDB?
	_, exists := myplaceutils.Users["laban"]

	if exists == false {
		t.Error("User not added to userDB after signed up")
	}
}



func TestPostMsg(t *testing.T) {
	myplaceutils.InitDBs()
	u1 := myplaceutils.AddNewUser("ask", "embla")
	u2 := myplaceutils.AddNewUser("adam", "eva")
	r1 := myplaceutils.AddNewRoom("livingroom")
	r2 := myplaceutils.AddNewRoom("bedroom")
	u1.JoinRoom(r1)
	u1.JoinRoom(r2)
	u2.JoinRoom(r1)

	str := "hello? who are you?"
	req := requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI := requests_responses.MsgInfo{0, r1.ID, u1.UName, -1, str}
	resp := requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse(t, req, resp)

	str = "anybody there?"
	req = requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{1, r1.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse(t, req, resp)

	str = "..."
	req = requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{2, r1.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse(t, req, resp)

	str = "no..yes"
	req = requests_responses.PostMsgRequest{12345, u2.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{3, r1.ID, u2.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse(t, req, resp)

	str = ""
	req = requests_responses.PostMsgRequest{12345, u1.UName, r2.ID, str}
	eresp := requests_responses.ErrorResponse{12345, requests_responses.PostMsgIndex, "bad msg length"}
	executeAndTestResponse(t, req, eresp)

	str = "que?\npor pue"
	req = requests_responses.PostMsgRequest{12345, u1.UName, r2.ID, str}
	msgI = requests_responses.MsgInfo{0, r2.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse(t, req, resp)

	str = "..."
	req = requests_responses.PostMsgRequest{12345, u2.UName, r2.ID, str}
	eresp = requests_responses.ErrorResponse{12345, requests_responses.PostMsgIndex, "user not in room"}
	executeAndTestResponse(t, req, eresp)
}



func TestJoinRoom(t *testing.T){
	myplaceutils.InitDBs()
	room := myplaceutils.AddNewRoom("213")	

	roomInfo := myplaceutils.CreateRoomInfo(room,nil,"Alex")
	req := requests_responses.JoinRoomRequest{12345,room.ID,"Alex"}
	eresp :=  requests_responses.ErrorResponse{12345,requests_responses.JoinRoomIndex,"There is no such user"}
	executeAndTestResponse(t,req,eresp)
	
	user1 := myplaceutils.AddNewUser("Eva", "1337")

	roomInfo = myplaceutils.CreateRoomInfo(room,nil,user1.UName)
	req = requests_responses.JoinRoomRequest{12345,room.ID,user1.UName}
	resp :=  requests_responses.JoinRoomResponse{12345,roomInfo,true}
	executeAndTestResponse(t,req,resp)
	
	user1.JoinRoom(room)
	
	roomInfo = myplaceutils.CreateRoomInfo(room,nil,user1.UName)
	req = requests_responses.JoinRoomRequest{12345,room.ID,user1.UName}
	resp = requests_responses.JoinRoomResponse{12345,roomInfo,true}
	executeAndTestResponse(t,req,resp)
	
	
	

}
