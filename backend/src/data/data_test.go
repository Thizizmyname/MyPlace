package data

import (
	"testing"
	"fmt"
	"time"
	"reflect"
	"container/list"
	"myplaceutils"
)


func TestLoad(t *testing.T) {
	//delete rooms and users
	us, rs, e := LoadDBs()

	if e == nil {
		fmt.Println("Remove 'rooms' and 'users'")
	}else if len(us) != 0 || len(rs) != 0 {
		t.Error("Wrong size after failed db-read")
	}
}

func TestStoreLoad(t *testing.T) {
	testStoreLoad(t, 100, 100, 100)
	testStoreLoad(t, 0, 0, 0)
	testStoreLoad(t, 1, 0, 0)
	testStoreLoad(t, 1, 100, 100)
	testStoreLoad(t, 100, 0, 100)
	testStoreLoad(t, 100, 100, 0)
}

func testStoreLoad(t *testing.T, no_users int, no_rooms int, no_msgs int) {
	var userDB myplaceutils.UserDB = make(myplaceutils.UserDB)
	var roomDB myplaceutils.RoomDB = make(myplaceutils.RoomDB)

	initUsers(userDB, 100, 10)
	initRooms(roomDB, 100, 3, 10)

	var err error
	if err = StoreDBs(userDB, roomDB); err != nil {
		panic(err)
	}

	userDB2, roomDB2, _ := LoadDBs()
	//roomDB2[5].Users.PushBack(100)

	if reflect.DeepEqual(userDB, userDB2) != true {
		t.Errorf("userDB != userDB2")
	}

	if roomDBsEqual(roomDB, roomDB2) != true {
		t.Errorf("roomDB != roomDB2")
	}
}

func initUsers(users myplaceutils.UserDB, no_users int, no_rooms int) {
	for i := 0; i < no_users; i++ {
		var u myplaceutils.User
		u.UName = fmt.Sprintf("user%v", i)
		u.Pass = fmt.Sprintf("pass%v", i)
		u.Rooms = list.New()

		for j := 0; j < no_rooms; j++ {
			u.Rooms.PushBack(j)
		}

		users[u.UName] = &u
	}
}

func initRooms(rooms myplaceutils.RoomDB, no_rooms, no_users int, no_msgs int) {
	for i := 0; i < no_rooms; i++ {
		var r myplaceutils.Room
		r.ID = i
		r.Name = fmt.Sprintf("room%v", i)

		r.Users = list.New()
		for j := 0; j < no_users; j++ {
			r.Users.PushBack(fmt.Sprintf("user%v", j))
		}

		r.Messages = make(map[int]*myplaceutils.Message)
		for j := 0; j < no_msgs; j++ {
			msg := myplaceutils.Message{
				j,
				time.Now(),
				fmt.Sprintf("user%v", j),
				"msg body" }

			r.Messages[j] = &msg
		}

		r.OutgoingChannels = list.New()

		rooms[i] = &r
	}
}


func equalUserList(u1 *list.List, u2 *list.List) bool {
	if u1.Len() != u2.Len() { return false }

	for e1, e2 := u1.Front(), u2.Front(); e1 != nil; e1, e2 = e1.Next(), e2.Next() {
		//uname1 := e1.Value.(string)
		//uname2 := e2.Value.(int)
		//fmt.Printf("%v, %v\n", uname1, uname2)

		//if uname1 != uname2 { return false }
	}

	return true
}

func roomsEqual(r1 *myplaceutils.Room, r2 *myplaceutils.Room) bool {
	return r1.ID == r2.ID &&
		r1.Name == r2.Name &&
		equalUserList(r1.Users, r2.Users) &&
		reflect.DeepEqual(r1.Messages, r2.Messages)
}

func roomDBsEqual(rs1 myplaceutils.RoomDB, rs2 myplaceutils.RoomDB) bool {
	if len(rs1) != len(rs2) { return false }

	for id, r1 := range rs1 {
		r2, ok := rs2[id]

		if !ok || !roomsEqual(r1, r2) { return false }
	}
	return true
}

func PrintList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func printUser(u *myplaceutils.User) {
	fmt.Println(u.UName)
	fmt.Println(u.Pass)
	PrintList(u.Rooms)
}

func printUsers(us myplaceutils.UserDB) {
	fmt.Println("userdb:")
	for _, u := range us {
		printUser(u)
	}
	fmt.Println("")
}
