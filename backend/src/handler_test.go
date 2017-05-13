package main

import (
	"testing"
	"os"
	"reflect"
	"myplaceutils"
	"requests_responses"
)

var handlerChan chan myplaceutils.HandlerArgs

func TestMain(m *testing.M) {
	myplaceutils.InitDBs()

	handlerChan = make(chan myplaceutils.HandlerArgs)
	go handler(handlerChan) //now handler is waiting for requests
	defer close(handlerChan)

	retCode := m.Run()
	os.Exit(retCode)
}


func executeAndTestResponse_chan(t *testing.T,
	responseChan chan requests_responses.Response,
	request requests_responses.Request,
	expectedResponse requests_responses.Response) {

	handlerArgs := myplaceutils.HandlerArgs{request, responseChan}

	handlerChan <- handlerArgs //send args to handler
	response := <-responseChan

	if r, ok := response.(requests_responses.PostMsgResponse); ok {
		req := request.(requests_responses.PostMsgRequest)
		expResp := expectedResponse.(requests_responses.PostMsgResponse)
		testPostMsgResponse(t, req, expResp, r, responseChan)
	} else if response != expectedResponse {
		t.Errorf("request: %v\nresponse: %v\nactual response: %v\nexpected response:%v",
			reflect.TypeOf(request),
			reflect.TypeOf(response),
			response,
			expectedResponse)
	}
}

func executeAndTestResponse(t *testing.T,
	request requests_responses.Request,
	expectedResponse requests_responses.Response) {

	responseChan := make(chan requests_responses.Response, 1)
	executeAndTestResponse_chan(t, responseChan, request, expectedResponse)

}

func testPostMsgResponse(t *testing.T, req requests_responses.PostMsgRequest, r2 requests_responses.PostMsgResponse, senderResp requests_responses.PostMsgResponse, senderResponseChan chan requests_responses.Response) {
	r := senderResp

	if r.RequestID != r2.RequestID ||
		r.Msg.MsgID != r2.Msg.MsgID ||
		r.Msg.RoomID != r2.Msg.RoomID ||
		r.Msg.UName != r2.Msg.UName ||
		r.Msg.Body != r2.Msg.Body {

		t.Errorf("request: %v\nresponse: %v\nactual response: %v\nexpected response:%v",
			reflect.TypeOf(r), reflect.TypeOf(r2), r, r2)
	}

	room := myplaceutils.Rooms[req.RoomID]

	for e := room.OutgoingChannels.Front(); e != nil; e = e.Next() {
		outChan := e.Value.(chan requests_responses.Response)

		if outChan == senderResponseChan {
			if responseChanIsEmpty(outChan) == false {
				t.Error("senderResponseChan not empty")
			}

		} else {
			r = (<-outChan).(requests_responses.PostMsgResponse)

			if r.RequestID != -1 ||
				r.Msg.MsgID != r2.Msg.MsgID ||
				r.Msg.RoomID != r2.Msg.RoomID ||
				r.Msg.UName != r2.Msg.UName ||
				r.Msg.Body != r2.Msg.Body {

				t.Errorf("request: %v\nresponse: %v\nactual response: %v\nexpected response:%v",
					reflect.TypeOf(r), reflect.TypeOf(r2), r, r2)
			}

			if responseChanIsEmpty(outChan) == false {
				t.Errorf("channel not empty")
			}
		}
	}
}


func responseChanIsEmpty(c chan requests_responses.Response) bool {
	select {
	case x := <-c:
		c <- x
		return false
	default:
		return true
	}
}







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
	if myplaceutils.UserExists("laban") == false {
		t.Error("User not added to userDB after signed up")
	}

	//test2
	req = requests_responses.SignUpRequest{1234, "alfons", "milla"}
	resp = requests_responses.SignUpResponse{1234, true, ""}
	executeAndTestResponse(t, req, resp)
	if myplaceutils.UserExists("alfons") == false {
		t.Error("User not added to userDB after signed up")
	}

	//test3
	req = requests_responses.SignUpRequest{123, "alfons", "mållgan"}
	resp = requests_responses.SignUpResponse{123, false, "uname"}
	executeAndTestResponse(t, req, resp)

	//test4
	req = requests_responses.SignUpRequest{123, "knyttet", ""}
	resp = requests_responses.SignUpResponse{123, false, "pass"}
	executeAndTestResponse(t, req, resp)

	if len(myplaceutils.Users) != 2 {
		t.Error("Wrong userDB length")
	}
}

func TestSignIn(t *testing.T) {
	myplaceutils.InitDBs()
	u1 := myplaceutils.AddNewUser("ask", "embla")
	u2 := myplaceutils.AddNewUser("adam", "eva")
	r1 := myplaceutils.AddNewRoom("livingroom")
	r2 := myplaceutils.AddNewRoom("bedroom")
	u1.JoinRoom(r1)
	u1.JoinRoom(r2)
	u2.JoinRoom(r1)

	if r1.OutgoingChannels.Len() != 0 || r2.OutgoingChannels.Len() != 0 {
		t.Error("Bad channels from start...?")
	}

	//test1
	req := requests_responses.SignInRequest{1234, u1.UName, u1.Pass}
	resp := requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse(t, req, resp)

	if r1.OutgoingChannels.Len() != 1 || r2.OutgoingChannels.Len() != 1 {
		t.Error("Bad outgoing channels after signin")
	}

	//test2
	req = requests_responses.SignInRequest{1234, u2.UName, u2.Pass}
	resp = requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse(t, req, resp)

	if r1.OutgoingChannels.Len() != 2 || r2.OutgoingChannels.Len() != 1 {
		t.Error("Bad outgoing channels after signin")
	}

	//test3 - wrong pass
	req = requests_responses.SignInRequest{1234, u1.UName, u2.Pass}
	resp = requests_responses.SignInResponse{1234, false, "pass"}
	executeAndTestResponse(t, req, resp)

	if r1.OutgoingChannels.Len() != 2 || r2.OutgoingChannels.Len() != 1 {
		t.Error("Bad outgoing channels after signin")
	}

	//test4 - nonexistent uname
	req = requests_responses.SignInRequest{1234, "krusmynta", u2.Pass}
	resp = requests_responses.SignInResponse{1234, false, "uname"}
	executeAndTestResponse(t, req, resp)

	if r1.OutgoingChannels.Len() != 2 || r2.OutgoingChannels.Len() != 1 {
		t.Error("Bad outgoing channels after signin")
	}
}

func TestCreateRoom(t *testing.T) {
	myplaceutils.InitDBs()
	u := myplaceutils.AddNewUser("ask", "embla")

	req := requests_responses.CreateRoomRequest{12345, "livingroom", u.UName}
	resp := requests_responses.CreateRoomResponse{12345, 0, "livingroom"}
	executeAndTestResponse(t, req, resp)

	req = requests_responses.CreateRoomRequest{12345, "bedroom", u.UName}
	resp = requests_responses.CreateRoomResponse{12345, 1, "bedroom"}
	executeAndTestResponse(t, req, resp)

	req = requests_responses.CreateRoomRequest{12345, "livingroom", "outsider_user"}
	eresp := requests_responses.ErrorResponse{12345, requests_responses.CreateRoomIndex, "no such user"}
	executeAndTestResponse(t, req, eresp)

	req = requests_responses.CreateRoomRequest{12345, "bedroom", u.UName}
	resp = requests_responses.CreateRoomResponse{12345, 2, "bedroom"}
	executeAndTestResponse(t, req, resp)
}

func TestPostMsg(t *testing.T) {
	myplaceutils.InitDBs()
	u1 := myplaceutils.AddNewUser("ask", "embla")
	u1_responseChan := make(chan requests_responses.Response, 1)
	u2 := myplaceutils.AddNewUser("adam", "eva")
	u2_responseChan := make(chan requests_responses.Response, 1)
	r1 := myplaceutils.AddNewRoom("livingroom")
	r2 := myplaceutils.AddNewRoom("bedroom")
	u1.JoinRoom(r1)
	u1.JoinRoom(r2)
	u2.JoinRoom(r1)

	//signin
	lrq := requests_responses.SignInRequest{1234, u1.UName, u1.Pass}
	lrp := requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u1_responseChan, lrq, lrp)
	lrq = requests_responses.SignInRequest{1234, u2.UName, u2.Pass}
	lrp = requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u2_responseChan, lrq, lrp)

	str := "hello? who are you?"
	req := requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI := requests_responses.MsgInfo{0, r1.ID, u1.UName, -1, str}
	resp := requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	str = "anybody there?"
	req = requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{1, r1.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	str = "..."
	req = requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{2, r1.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	str = "no..yes"
	req = requests_responses.PostMsgRequest{12345, u2.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{3, r1.ID, u2.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u2_responseChan, req, resp)

	str = ""
	req = requests_responses.PostMsgRequest{12345, u1.UName, r2.ID, str}
	eresp := requests_responses.ErrorResponse{12345, requests_responses.PostMsgIndex, "bad msg length"}
	executeAndTestResponse_chan(t, u1_responseChan, req, eresp)

	str = "que?\npor pue"
	req = requests_responses.PostMsgRequest{12345, u1.UName, r2.ID, str}
	msgI = requests_responses.MsgInfo{0, r2.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	str = "..."
	req = requests_responses.PostMsgRequest{12345, u2.UName, r2.ID, str}
	eresp = requests_responses.ErrorResponse{12345, requests_responses.PostMsgIndex, "user not in room"}
	executeAndTestResponse_chan(t, u2_responseChan, req, eresp)
}

func TestSignOut(t *testing.T) {
	myplaceutils.InitDBs()
	u1 := myplaceutils.AddNewUser("ask", "embla")
	u1_responseChan := make(chan requests_responses.Response, 1)
	u2 := myplaceutils.AddNewUser("adam", "eva")
	u2_responseChan := make(chan requests_responses.Response, 1)
	r1 := myplaceutils.AddNewRoom("livingroom")
	r2 := myplaceutils.AddNewRoom("bedroom")
	u1.JoinRoom(r1)
	u1.JoinRoom(r2)
	u2.JoinRoom(r1)

	//signin
	lrq := requests_responses.SignInRequest{1234, u1.UName, u1.Pass}
	lrp := requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u1_responseChan, lrq, lrp)
	lrq = requests_responses.SignInRequest{1234, u2.UName, u2.Pass}
	lrp = requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u2_responseChan, lrq, lrp)
	if r1.OutgoingChannels.Len() != 2 || r2.OutgoingChannels.Len() != 1 {
		t.Error()
	}

	//signout
	req := requests_responses.SignOutRequest{12345, u1.UName}
	resp := requests_responses.SignOutResponse{12345}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)
	if r1.OutgoingChannels.Len() != 1 || r2.OutgoingChannels.Len() != 0 {
		t.Errorf("Bad channel count, r1=%v, r2=%v",
			r1.OutgoingChannels.Len(),
			r2.OutgoingChannels.Len())
	}

	req = requests_responses.SignOutRequest{12345, u2.UName}
	resp = requests_responses.SignOutResponse{12345}
	executeAndTestResponse_chan(t, u2_responseChan, req, resp)
	if r1.OutgoingChannels.Len() != 0 || r2.OutgoingChannels.Len() != 0 {
		t.Errorf("Bad channel count, r1=%v, r2=%v",
			r1.OutgoingChannels.Len(),
			r2.OutgoingChannels.Len())
	}

	req = requests_responses.SignOutRequest{12345, u2.UName}
	resp = requests_responses.SignOutResponse{12345}
	executeAndTestResponse_chan(t, u2_responseChan, req, resp)
	if r1.OutgoingChannels.Len() != 0 || r2.OutgoingChannels.Len() != 0 {
		t.Errorf("Bad channel count, r1=%v, r2=%v",
			r1.OutgoingChannels.Len(),
			r2.OutgoingChannels.Len())
	}
}
