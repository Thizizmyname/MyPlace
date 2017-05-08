package main

import (
	"testing"
	"reflect"
	"myplaceutils"
	"requests_responses"
)

func testResponse(t *testing.T, request requests_responses.Request, expectedResponse requests_responses.Response) {
	handlerChan := make(chan myplaceutils.HandlerArgs)
	go handler(handlerChan) //now handler is waiting for requests
	defer close(handlerChan)

	responseChan := make(chan requests_responses.Response)
	handlerArgs := myplaceutils.HandlerArgs{request, responseChan}

	handlerChan <- handlerArgs //send args to handler
	response := <-responseChan

	if response != expectedResponse {
		t.Error(reflect.TypeOf(request))
	}
}


/*
00.signup
args: uname, pass
response: -
note: error if uname is taken/ pass to short/ illegal characters/ ...
side-effect: updates users_db
*/
func TestSignUp(t *testing.T){
	req := requests_responses.SignUpRequest{12345, "user", "pass"}
	resp := requests_responses.SignUpResponse{12345, true, ""}
	testResponse(t, req, resp)
}
