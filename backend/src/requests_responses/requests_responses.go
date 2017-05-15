package requests_responses

import (
  "encoding/json"
  "strconv"
  "fmt"
  "errors"
)

//----------------------INTERFACE START

type SignUpRequest struct {
RequestID int
UName string
Pass string
}

type SignUpResponse struct {
RequestID int
Result bool
ErrorCause string
}

type SignInRequest struct {
RequestID int
UName string
Pass string
}

type SignInResponse struct {
RequestID int
Result bool
ErrorCause string
}

type GetRoomsRequest struct {
RequestID int
UName string
}

type GetRoomsResponse struct {
	RequestID int
	Rooms []RoomInfo
}

type GetRoomUsersRequest struct {
RequestID int
RoomID int
}

type GetRoomUsersResponse struct {
  RequestID int
  RoomID int
  UNames []string
}

type GetOlderMsgsRequest struct {
  RequestID int
  RoomID int
  MsgID int
}

type GetOlderMsgsResponse struct {
	RequestID int
	Messages []MsgInfo
}

type GetNewerMsgsRequest struct {
  RequestID int
  RoomID int
  MsgID int
}

type GetNewerMsgsResponse struct {
	RequestID int
	Messages []MsgInfo
}

type JoinRoomRequest struct {
  RequestID int
  RoomID int
  UName string
}

type JoinRoomResponse struct {
	RequestID int
	JoinedRoom RoomInfo
	RoomIDAccepted bool
}

type LeaveRoomRequest struct {
  RequestID int
  RoomID int
  UName string
}

type LeaveRoomResponse struct {
  RequestID int
}

type CreateRoomRequest struct {
  RequestID int
  RoomName string
  UName string
}

type CreateRoomResponse struct {
  RequestID int
  RoomID int
  RoomName string
}

type PostMsgRequest struct {
  RequestID int
  UName string
  RoomID int
  Body string
}

type PostMsgResponse struct {
	RequestID int //-1 if no request was made
	Msg MsgInfo
}

type MsgReadRequest struct {
  RequestID int
  MsgID int
  RoomID int
  UName string
}

type MsgReadResponse struct {
  RequestID int
}

type SignOutRequest struct {
  RequestID int
  UName string
}

type SignOutResponse struct {
  RequestID int
}

type ErrorResponse struct {
  RequestID int
  RequestIndex int
  ErrorCause string
}

const (
	SignUpIndex = 0
	SignInIndex = 1
	GetRoomsIndex = 2
	GetRoomUsersIndex = 3
	GetOlderMsgsIndex = 4
	GetNewerMsgsIndex = 5
	JoinRoomIndex = 6
	LeaveRoomIndex = 7
	CreateRoomIndex = 8
	PostMsgIndex = 9
	MsgReadIndex = 10
	SignOutIndex = 11
	ErrorIndex = -1
)

//---------------------------INTERFACE STOP

type RoomInfo struct {
	ID int
	Name string
	LatestMsg MsgInfo
	LatestReadMsgID int
}

type MsgInfo struct {
	MsgID int
	RoomID int
	UName string
	Time int64 //number of milliseconds since January 1, 1970 UTC
	Body string
}

type Request interface {}
type Response interface {}

func FromRequestString(requestString string) (Request, error) {
	reqType, err := strconv.Atoi(requestString[:2])

	if (err != nil) {
		return nil, err
	}

	jsonRequest := requestString[2 : len(requestString) - 1]

	switch reqType {
	case 0:
		var r SignUpRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 1:
		var r SignInRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 2:
		var r GetRoomsRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 3:
		var r GetRoomUsersRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 4:
		var r GetOlderMsgsRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 5:
		var r GetNewerMsgsRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 6:
		var r JoinRoomRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 7:
		var r LeaveRoomRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 8:
		var r CreateRoomRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 9:
		var r PostMsgRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 10:
		var r MsgReadRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	case 11:
		var r SignOutRequest
		err := json.Unmarshal([]byte(jsonRequest), &r)
		return r, err
	default:
		return nil, errors.New("illegal reqType")
	}
}

func ToResponseString(response Response) (string, error) {
  jsonResponse, err := json.Marshal(response)

  if (err != nil) {
    return "", err
  }

  switch response.(type) {
  case SignUpResponse:
    return fmt.Sprintf("00%s\n", jsonResponse), nil
  case SignInResponse:
    return fmt.Sprintf("01%s\n", jsonResponse), nil
  case GetRoomsResponse:
    return fmt.Sprintf("02%s\n", jsonResponse), nil
  case GetRoomUsersResponse:
    return fmt.Sprintf("03%s\n", jsonResponse), nil
  case GetOlderMsgsResponse:
    return fmt.Sprintf("04%s\n", jsonResponse), nil
  case GetNewerMsgsResponse:
    return fmt.Sprintf("05%s\n", jsonResponse), nil
  case JoinRoomResponse:
    return fmt.Sprintf("06%s\n", jsonResponse), nil
  case LeaveRoomResponse:
    return fmt.Sprintf("07%s\n", jsonResponse), nil
  case CreateRoomResponse:
    return fmt.Sprintf("08%s\n", jsonResponse), nil
  case PostMsgResponse:
    return fmt.Sprintf("09%s\n", jsonResponse), nil
  case MsgReadResponse:
    return fmt.Sprintf("10%s\n", jsonResponse), nil
  case SignOutResponse:
    return fmt.Sprintf("11%s\n", jsonResponse), nil
  case ErrorResponse:
    return fmt.Sprintf("-1%s\n", jsonResponse), nil
  default:
    return "", errors.New("illegal response type")
  }
}

func FromResponseString(responseString string) (Response, error) {
  respType, err := strconv.Atoi(responseString[:2])

  if (err != nil) {
    return nil, err
  }

  jsonResponse := responseString[2:]

  switch respType {
  case 0:
    var r SignUpResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 1:
    var r SignInResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 2:
    var r GetRoomsResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 3:
    var r GetRoomUsersResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 4:
    var r GetOlderMsgsResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 5:
    var r GetNewerMsgsResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 6:
    var r JoinRoomResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 7:
    var r LeaveRoomResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 8:
    var r CreateRoomResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 9:
    var r PostMsgResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 10:
    var r MsgReadResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case 11:
    var r SignOutResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  case -1:
    var r ErrorResponse
    err := json.Unmarshal([]byte(jsonResponse), &r)
    return r, err
  default:
    return nil, errors.New("illegal reqType")
  }
}

func ToRequestString(request Request) (string, error) {
  jsonRequest, err := json.Marshal(request)

  if (err != nil) {
    return "", err
  }

  switch request.(type) {
  case SignUpRequest:
    return fmt.Sprintf("00%s", jsonRequest), nil
  case SignInRequest:
    return fmt.Sprintf("01%s", jsonRequest), nil
  case GetRoomsRequest:
    return fmt.Sprintf("02%s", jsonRequest), nil
  case GetRoomUsersRequest:
    return fmt.Sprintf("03%s", jsonRequest), nil
  case GetOlderMsgsRequest:
    return fmt.Sprintf("04%s", jsonRequest), nil
  case GetNewerMsgsRequest:
    return fmt.Sprintf("05%s", jsonRequest), nil
  case JoinRoomRequest:
    return fmt.Sprintf("06%s", jsonRequest), nil
  case LeaveRoomRequest:
    return fmt.Sprintf("07%s", jsonRequest), nil
  case CreateRoomRequest:
    return fmt.Sprintf("08%s", jsonRequest), nil
  case PostMsgRequest:
    return fmt.Sprintf("09%s", jsonRequest), nil
  case MsgReadRequest:
    return fmt.Sprintf("10%s", jsonRequest), nil
  case SignOutRequest:
    return fmt.Sprintf("11%s", jsonRequest), nil
  default:
    return "", errors.New("illegal request type")
  }
}
