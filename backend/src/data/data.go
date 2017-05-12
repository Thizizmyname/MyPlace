package data

import (
	"encoding/json"
	"io/ioutil"
	"myplaceutils"
	"container/list"
)

type UserStore struct {
  UName string
  Pass string
  Rooms []int
}

type RoomStore struct {
  ID int
  Name     string
  Users []string
  Messages map[int]*myplaceutils.Message //ID is key
}

type UserDBStore map[string]*UserStore
type RoomDBStore map[int]*RoomStore


//Stores dbs to files named "rooms" and "users" in same folder as this.
//Encrypted using json.
func StoreDBs(us myplaceutils.UserDB, rs myplaceutils.RoomDB) error {
	usStore, rsStore := toStoreFormat(us, rs)

	var usj []byte
	var rsj []byte
	var e error

	if usj, e = json.Marshal(usStore); e != nil {
		return e
	}

	if rsj, e = json.Marshal(rsStore); e != nil {
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


//Reads from files "users" and "rooms" is same folder as this file.
//If read or decryption fails, empty dbs are returned.
func LoadDBs() (myplaceutils.UserDB, myplaceutils.RoomDB, error) {
	var usj []byte
	var rsj []byte
	var e error
	us := make(myplaceutils.UserDB)
	rs := make(myplaceutils.RoomDB)

	if usj, e = ioutil.ReadFile("users"); e != nil {
		return us, rs, e
	}

	if rsj, e = ioutil.ReadFile("rooms"); e != nil {
		return us, rs, e
	}

	var usStore UserDBStore
	var rsStore RoomDBStore

	if e = json.Unmarshal(usj, &usStore); e != nil {
		return us, rs, e
	}

	if e = json.Unmarshal(rsj, &rsStore); e != nil {
		return us, rs, e
	}

	us, rs = fromStoreFormat(usStore, rsStore)
	return us, rs, nil
}

func toStoreFormat(us myplaceutils.UserDB, rs myplaceutils.RoomDB) (UserDBStore, RoomDBStore) {
	usStore := make(UserDBStore)
	rsStore := make(RoomDBStore)

	for uname, user := range us {
		rs := toRoomSlice(user.Rooms)
		uStore := UserStore{uname, user.Pass, rs}
		usStore[uname] = &uStore
	}

	for roomID, room := range rs {
		us := toUserSlice(room.Users)
		rStore := RoomStore{roomID, room.Name, us, room.Messages}
		rsStore[roomID] = &rStore
	}

	return usStore, rsStore
}

func fromStoreFormat(usStore UserDBStore, rsStore RoomDBStore) (myplaceutils.UserDB, myplaceutils.RoomDB) {
	us := make(myplaceutils.UserDB)
	rs := make(myplaceutils.RoomDB)

	for uname, uStore := range usStore {
		rs := fromRoomSlice(uStore.Rooms)
		u := myplaceutils.User{uname, uStore.Pass, rs}
		us[uname] = &u
	}

	for roomID, rStore := range rsStore {
		us := fromUserSlice(rStore.Users)
		r := myplaceutils.Room{roomID, rStore.Name, us, rStore.Messages, list.New()}
		rs[roomID] = &r
	}

	return us, rs
}

func toRoomSlice(rList *list.List) []int {
	var rSlice []int
	for e := rList.Front(); e != nil; e = e.Next() {
		rSlice = append(rSlice, e.Value.(int))
	}
	return rSlice
}

func toUserSlice(uList *list.List) []string {
	var uSlice []string
	for e := uList.Front(); e != nil; e = e.Next() {
		uSlice = append(uSlice, e.Value.(string))
	}
	return uSlice
}

func fromRoomSlice(rSlice []int) *list.List {
	rList := list.New()
	for id := range rSlice {
		rList.PushBack(id)
	}
	return rList
}

func fromUserSlice(uSlice []string) *list.List {
	uList := list.New()

	for _, uname := range uSlice {
		uList.PushBack(uname)
	}

	return uList
}
