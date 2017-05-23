package handler

import (
	"testing"
	"os"
	"reflect"
	"myplaceutils"
	"requests_responses"
	"io/ioutil"
	"io"
	"log"
)

var handlerChan chan myplaceutils.HandlerArgs

func TestMain(m *testing.M) {
	myplaceutils.InitDBs()
	initLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	handlerChan = make(chan myplaceutils.HandlerArgs)
	go ResponseHandler(handlerChan) //now handler is waiting for requests
	defer close(handlerChan)

	retCode := m.Run()
	os.Exit(retCode)
}

//Initialize loggers
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
	u1_responseChan := make(chan requests_responses.Response, 1)
	u2 := myplaceutils.AddNewUser("adam", "eva")
	r1 := myplaceutils.AddNewRoom("livingroom")
	r2 := myplaceutils.AddNewRoom("bedroom")
	u1.JoinRoom(r1)
	u1.JoinRoom(r2)
	u2.JoinRoom(r1)

	if r1.OutgoingChannels.Len() != 0 || r2.OutgoingChannels.Len() != 0 {
		t.Error("Bad channels from start")
	}

	//test1
	req := requests_responses.SignInRequest{1234, u1.UName, u1.Pass}
	resp := requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

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

	//existing uname (re-signin)
	req = requests_responses.SignInRequest{1234, u1.UName, u1.Pass}
	resp = requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	if r1.OutgoingChannels.Len() != 2 || r2.OutgoingChannels.Len() != 1 {
		t.Errorf("Bad outgoing channels after signin, %v, %v", r1.OutgoingChannels.Len(), r2.OutgoingChannels.Len())
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

func TestPostMsgAndMsgRead(t *testing.T) {
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

	//post msgs
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




	//test msgRead
	//r1 has 4 msgs, r2 has 1
	//u1 in r1&r2, u2 in r1
	//u1,r1: latestRead=3
	mrreq := requests_responses.MsgReadRequest{1, 3, r1.ID, u1.UName}
	mrresp := requests_responses.MsgReadResponse{1}
	executeAndTestResponse(t, mrreq, mrresp)
	if u1.GetLatestReadMsg(r1.ID) != 3 {
		t.Error("bad latestReadMsg")
	}

	//same
	mrreq = requests_responses.MsgReadRequest{2, 3, r1.ID, u1.UName}
	mrresp = requests_responses.MsgReadResponse{2}
	executeAndTestResponse(t, mrreq, mrresp)
	if u1.GetLatestReadMsg(r1.ID) != 3 {
		t.Error("bad latestReadMsg")
	}

	//u1,r1: latestRead=0
	mrreq = requests_responses.MsgReadRequest{3, 0, r1.ID, u1.UName}
	mrresp = requests_responses.MsgReadResponse{3}
	executeAndTestResponse(t, mrreq, mrresp)
	if u1.GetLatestReadMsg(r1.ID) != 0 {
		t.Error("bad latestReadMsg")
	}

	//u1,r1: latestRead=4
	mrreq = requests_responses.MsgReadRequest{4, 4, r1.ID, u1.UName}
	eresp = requests_responses.ErrorResponse{4, requests_responses.MsgReadIndex, "bad msgID or user not in room"}
	executeAndTestResponse(t, mrreq, eresp)
	if u1.GetLatestReadMsg(r1.ID) != 0 {
		t.Error("bad latestReadMsg")
	}

	//u2,r1: latestRead=10
	mrreq = requests_responses.MsgReadRequest{5, 10, r1.ID, u2.UName}
	eresp = requests_responses.ErrorResponse{5, requests_responses.MsgReadIndex, "bad msgID or user not in room"}
	executeAndTestResponse(t, mrreq, eresp)
	if u2.GetLatestReadMsg(r1.ID) != -1 {
		t.Error("bad latestReadMsg")
	}

	//u2,r2: latestRead=0
	mrreq = requests_responses.MsgReadRequest{6, 0, r2.ID, u2.UName}
	eresp = requests_responses.ErrorResponse{6, requests_responses.MsgReadIndex, "bad msgID or user not in room"}
	executeAndTestResponse(t, mrreq, eresp)
	if u2.GetLatestReadMsg(r1.ID) != -1 {
		t.Error("bad latestReadMsg")
	}

	//u2,r2: latestRead=10
	mrreq = requests_responses.MsgReadRequest{12345, 10, r2.ID, u2.UName}
	eresp = requests_responses.ErrorResponse{12345, requests_responses.MsgReadIndex, "bad msgID or user not in room"}
	executeAndTestResponse(t, mrreq, eresp)
	if u2.GetLatestReadMsg(r1.ID) != -1 {
		t.Error("bad latestReadMsg")
	}
}

func TestMsgRead(t *testing.T) {

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



func TestJoinRoom(t *testing.T){
	myplaceutils.InitDBs()

	//Test if the user doesn't exist
	req := requests_responses.JoinRoomRequest{12345,0,"Alex"}
	eresp :=  requests_responses.ErrorResponse{12345,requests_responses.JoinRoomIndex,"There is no such user"}
	executeAndTestResponse(t,req,eresp)

	user1 := myplaceutils.AddNewUser("Eva", "1337")

	// Test if the room doesn't exist
	roomInfo := requests_responses.RoomInfo{}
	req = requests_responses.JoinRoomRequest{12345,1,user1.UName}
	resp := requests_responses.JoinRoomResponse{12345,roomInfo,false}
	executeAndTestResponse(t,req,resp)
	
	room0 := myplaceutils.AddNewRoom("213")
	
	//Test if Eva joins Room 213
	roomInfo = myplaceutils.CreateRoomInfo(room0,user1)
	req = requests_responses.JoinRoomRequest{12345,room0.ID,user1.UName}
	resp =  requests_responses.JoinRoomResponse{12345,roomInfo,true}
	executeAndTestResponse(t,req,resp)

	room1 := myplaceutils.AddNewRoom("2001")

	// Test if a user can join another room aswell
	roomInfo = myplaceutils.CreateRoomInfo(room1,user1)
	req = requests_responses.JoinRoomRequest{12345,room1.ID,user1.UName}
	resp = requests_responses.JoinRoomResponse{12345,roomInfo,true}
	executeAndTestResponse(t,req,resp)

	// Test if a user joins a room which his is allready a member of
	roomInfo = myplaceutils.CreateRoomInfo(room0,user1)
	req = requests_responses.JoinRoomRequest{12345,room0.ID,user1.UName}
	eresp = requests_responses.ErrorResponse{12345,requests_responses.JoinRoomIndex,"User is already a member of the room"}
	executeAndTestResponse(t,req,eresp)
	
	
	
}

func TestLeaveRoom(t *testing.T){
	myplaceutils.InitDBs()


	// Test 1 Error -The user doesn't exist
	req := requests_responses.LeaveRoomRequest{12345,0,"Alex"}
	eresp := requests_responses.ErrorResponse{12345,requests_responses.LeaveRoomIndex,"There is no such user"}
	executeAndTestResponse(t,req,eresp)

	u0 := myplaceutils.AddNewUser("Alex","qwerty")
	u1 := myplaceutils.AddNewUser("Erik", "1337")

	// Test if the room doesn't exist
	req = requests_responses.LeaveRoomRequest{12345,1,"Alex"}
	eresp = requests_responses.ErrorResponse{12345,requests_responses.LeaveRoomIndex,"Bad roomID"}
	executeAndTestResponse(t,req,eresp)

	room := myplaceutils.AddNewRoom("213")


	// Test 2 Error - The user isn't in the room
	req = requests_responses.LeaveRoomRequest{12345,room.ID,u0.UName}
	eresp = requests_responses.ErrorResponse{12345, requests_responses.LeaveRoomIndex,"There is no such user in the room"}
	executeAndTestResponse(t,req,eresp)

	u0.JoinRoom(room)
	u1.JoinRoom(room)

	// Test 3 - Checks if it can remove the user
	req = requests_responses.LeaveRoomRequest{12345,room.ID,u0.UName}
	resp := requests_responses.LeaveRoomResponse{12345}
	executeAndTestResponse(t,req,resp)
}


func TestGetNewerMsgs(t *testing.T){
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

//	var msgInfos = make([] requests_responses.MsgInfo,10) // 10 = antal meddelanden
	
	//signin
	lrq := requests_responses.SignInRequest{1234, u1.UName, u1.Pass}
	lrp := requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u1_responseChan, lrq, lrp)
	lrq = requests_responses.SignInRequest{1234, u2.UName, u2.Pass}
	lrp = requests_responses.SignInResponse{1234, true, ""}
	executeAndTestResponse_chan(t, u2_responseChan, lrq, lrp)

	//post msgs
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

	str = "Yea, what's your name?"
	req = requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{2, r1.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	str = "My name is Alex"
	req = requests_responses.PostMsgRequest{12345, u1.UName, r1.ID, str}
	msgI = requests_responses.MsgInfo{3, r1.ID, u1.UName, -1, str}
	resp = requests_responses.PostMsgResponse{12345, msgI}
	executeAndTestResponse_chan(t, u1_responseChan, req, resp)

	/*
	req1 := requests_responses.GetNewerMsgsRequest{12345,r1.ID,2} //Tror vi har bytt namn from room.ID till room.roomID 
	resp1 := requests_responses.GetNewerMsgsResponse{12345,msgInfos}
	executeAndTestResponse(t, req1,resp1)
*/
} 


