package data

import (
	"encoding/json"
	"io/ioutil"
	"myplaceutils"
)

func StoreDBs(us myplaceutils.UserDB, rs myplaceutils.RoomDB) error {
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

func LoadDBs() (myplaceutils.UserDB, myplaceutils.RoomDB, error) {
	var usj []byte
	var rsj []byte
	var e error

	if usj, e = ioutil.ReadFile("users"); e != nil {
		return nil, nil, e
	}

	if rsj, e = ioutil.ReadFile("rooms"); e != nil {
		return nil, nil, e
	}

	var us myplaceutils.UserDB
	var rs myplaceutils.RoomDB

	if e = json.Unmarshal(usj, &us); e != nil {
		return nil, nil, e
	}

	if e = json.Unmarshal(rsj, &rs); e != nil {
		return nil, nil, e
	}

	return us, rs, nil
}
