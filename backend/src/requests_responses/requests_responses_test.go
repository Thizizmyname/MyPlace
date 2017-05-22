package requests_responses

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestFromRequestString(t *testing.T) {
	//signup
	signUp1 := SignUpRequest{0, "user", "pass"}
	jsonReq, _ := json.Marshal(signUp1)
	reqString := fmt.Sprintf("00%s\n", jsonReq)
	signUp2, _ := FromRequestString(reqString)

	if (signUp1 != signUp2.(SignUpRequest)) {
		t.Errorf("signup")
	}

	//signin
	signIn1 := SignInRequest{0, "user", "pass"}
	jsonReq, _ = json.Marshal(signIn1)
	reqString = fmt.Sprintf("01%s\n", jsonReq)
	signIn2, _ := FromRequestString(reqString)

	if (signIn1 != signIn2.(SignInRequest)) {
		t.Errorf("signin")
	}

	//getRooms
	getRooms1 := GetRoomsRequest{0, "user"}
	jsonReq, _ = json.Marshal(getRooms1)
	reqString = fmt.Sprintf("02%s\n", jsonReq)
	getRooms2, _ := FromRequestString(reqString)

	if (getRooms1 != getRooms2.(GetRoomsRequest)) {
		t.Errorf("getRooms")
	}

	//getRoomUsers
	getRoomUsers1 := GetRoomUsersRequest{12345, 54321}
	jsonReq, _ = json.Marshal(getRoomUsers1)
	reqString = fmt.Sprintf("03%s\n", jsonReq)
	getRoomUsers2, _ := FromRequestString(reqString)

	if (getRoomUsers1 != getRoomUsers2.(GetRoomUsersRequest)) {
		t.Errorf("getRoomUsers")
	}

	//getOlderMsgs
	getOlderMsgs1 := GetOlderMsgsRequest{12345, 54321, 987654}
	jsonReq, _ = json.Marshal(getOlderMsgs1)
	reqString = fmt.Sprintf("04%s\n", jsonReq)
	getOlderMsgs2, _ := FromRequestString(reqString)

	if (getOlderMsgs1 != getOlderMsgs2.(GetOlderMsgsRequest)) {
		t.Errorf("getOlderMsgs")
	}

	//getNewerMsgs
	getNewerMsgs1 := GetNewerMsgsRequest{12345, 54321, 987654}
	jsonReq, _ = json.Marshal(getNewerMsgs1)
	reqString = fmt.Sprintf("05%s\n", jsonReq)
	getNewerMsgs2, _ := FromRequestString(reqString)

	if (getNewerMsgs1 != getNewerMsgs2.(GetNewerMsgsRequest)) {
		t.Errorf("getNewerMsgs")
	}

	//joinRoom
	joinRoom1 := JoinRoomRequest{12345, 54321, "user"}
	jsonReq, _ = json.Marshal(joinRoom1)
	reqString = fmt.Sprintf("06%s\n", jsonReq)
	joinRoom2, _ := FromRequestString(reqString)

	if (joinRoom1 != joinRoom2.(JoinRoomRequest)) {
		t.Errorf("joinRoom")
	}

	//leaveRoom
	leaveRoom1 := LeaveRoomRequest{12345, 54321, "user"}
	jsonReq, _ = json.Marshal(leaveRoom1)
	reqString = fmt.Sprintf("07%s\n", jsonReq)
	leaveRoom2, _ := FromRequestString(reqString)

	if (leaveRoom1 != leaveRoom2.(LeaveRoomRequest)) {
		t.Errorf("leaveRoom")
	}

	//createRoom
	createRoom1 := CreateRoomRequest{0, "room", "user"}
	jsonReq, _ = json.Marshal(createRoom1)
	reqString = fmt.Sprintf("08%s\n", jsonReq)
	createRoom2, _ := FromRequestString(reqString)

	if (createRoom1 != createRoom2.(CreateRoomRequest)) {
		t.Errorf("createRoom")
	}

	//postMsg
	postMsg1 := PostMsgRequest{12345, "user", 54321, "zsp myplace?"}
	jsonReq, _ = json.Marshal(postMsg1)
	reqString = fmt.Sprintf("09%s\n", jsonReq)
	postMsg2, _ := FromRequestString(reqString)

	if (postMsg1 != postMsg2.(PostMsgRequest)) {
		t.Errorf("postMsg")
	}

	//msgRead
	msgRead1 := MsgReadRequest{12345, 54321, 0, "user"}
	jsonReq, _ = json.Marshal(msgRead1)
	reqString = fmt.Sprintf("10%s\n", jsonReq)
	msgRead2, _ := FromRequestString(reqString)

	if (msgRead1 != msgRead2.(MsgReadRequest)) {
		t.Errorf("msgRead")
	}

	// //index out of bounds
	// reqString = fmt.Sprintf("11%s", jsonReq)
	// _, err := FromRequestString(reqString)

	// if (err == nil) {
	// 	t.Errorf("index out of bounds")
	// }

	// //incorrect index
	// reqString = fmt.Sprintf("00%s", jsonReq)
	// _, err = FromRequestString(reqString)

	// if (err == nil) {
	// 	t.Errorf("incorrect index")
	// }

	//bad formatting
	reqString = fmt.Sprintf("1%s\n", jsonReq)
	_, err := FromRequestString(reqString)

	if (err == nil) {
		t.Errorf("bad formatting")
	}
}

// For printing out encrypted requests and responses:
//
// func TestPrintouts(t *testing.T) {
// 	var msgI0 msgInfo = msgInfo{0, 0, "user0", 123456789, "msg body"}
// 	var msgI1 msgInfo = msgInfo{1, 1, "user1", 123456789, "msg body"}
// 	var msgI2 msgInfo = msgInfo{2, 2, "user2", 123456789, "msg body"}
// 	var roomI0 roomInfo = roomInfo{0, "room1", msgI0, 0}
// 	var roomI1 roomInfo = roomInfo{1, "room1", msgI1, 0}
// 	var roomI2 roomInfo = roomInfo{2, "room2", msgI2, 0}
// 	var req string
// 	var resp string

// 	fmt.Println("SignUp")
// 	signUp := SignUpRequest{12345, "user", "pass"}
// 	jsonReq, _ := json.Marshal(signUp)
// 	req = fmt.Sprintf("00%s", jsonReq)
// 	resp, _ = ToResponseString(SignUpResponse{12345, true, ""})
// 	fmt.Println(string(req))
// 	fmt.Println(resp)

// 	fmt.Println("SignIn")
// 	signIn := SignInRequest{12345, "user", "pass"}
// 	jsonReq, _ = json.Marshal(signIn)
// 	req = fmt.Sprintf("01%s", jsonReq)
// 	resp, _ = ToResponseString(SignInResponse{12345, false, "user"})
// 	fmt.Println(req)
// 	fmt.Println(resp)

// 	fmt.Println("GetRooms")
// 	getRooms := GetRoomsRequest{12345, "user"}
// 	jsonReq, _ = json.Marshal(getRooms)
// 	req = fmt.Sprintf("02%s", jsonReq)
// 	resp, _ = ToResponseString(GetRoomsResponse{12345, []roomInfo{roomI0, roomI1, roomI2}})
// 	fmt.Println(string(req))
// 	fmt.Println(resp)

// 	fmt.Println("GetRoomUsers")
// 	getRoomUsers := GetRoomUsersRequest{12345, 54321}
// 	jsonReq, _ = json.Marshal(getRoomUsers)
// 	req = fmt.Sprintf("03%s", jsonReq)
// 	resp, _ = ToResponseString(GetRoomUsersResponse{12345, 54321, []string{"user0", "user1", "user2"}})
// 	fmt.Println(string(req))
// 	fmt.Println(resp)

// 	fmt.Println("GetOlderMsgsRequest")
// 	getOlderMsgs := GetOlderMsgsRequest{12345, 54321, 987654}
// 	jsonReq, _ = json.Marshal(getOlderMsgs)
// 	req = fmt.Sprintf("04%s", jsonReq)
// 	resp, _ = ToResponseString(GetOlderMsgsResponse{12345, []msgInfo{msgI0, msgI1, msgI2}})
// 	fmt.Println(string(req))
// 	fmt.Println(resp)
// }
