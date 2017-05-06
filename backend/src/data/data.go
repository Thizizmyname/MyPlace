package data

import (
	"encoding/json"
	"io/ioutil"
  "net"
  "fmt"
 // "myplaceutils"
)
type UserDB []myplaceutils.User
type RoomDB []myplaceutils.Room

/*
//DETTA ÄR FLYTTAT TILL src/myplaceutils/myplaceutils.go för att modularisera

type User struct {
	Uname string
	Pass string
	Rooms []Room
  ActiveConn net.Conn
}

type Room struct {
	Name string
	NoPeople int
  Users []User
	Messages []Message
}

type Message struct {
	Time time.Time
	Uname string
	Body string
  ID string
  
}
*/


func StoreDBs(us UserDB, rs RoomDB) error {
	var usj []byte
	var rsj []byte
	var e error

	if usj, e = json.Marshal(us); e != nil {
		return e
	}

	if rsj, e = json.Marshal(rs); e != nil {
		return e
	}

	if e = ioutil.WriteFile("users", usj, 0644); e != nil {
		return e
	}

	if e = ioutil.WriteFile("rooms", rsj, 0644); e != nil {
		return e
	}

	return nil

}

func LoadDBs() (UserDB, RoomDB, error) {
	var usj []byte
	var rsj []byte
	var e error

	if usj, e = ioutil.ReadFile("users"); e != nil {
		return nil, nil, e
	}

	if rsj, e = ioutil.ReadFile("rooms"); e != nil {
		return nil, nil, e
	}

	var us UserDB
	var rs RoomDB

	if e = json.Unmarshal(usj, &us); e != nil {
		return nil, nil, e
	}

	if e = json.Unmarshal(rsj, &rs); e != nil {
		return nil, nil, e
	}

	return us, rs, nil
}

