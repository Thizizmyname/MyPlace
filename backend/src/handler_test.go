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
