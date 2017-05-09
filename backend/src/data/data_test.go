package data

import (
	"testing"
	"fmt"
	"time"
	//"reflect"
	"container/list"
	"myplaceutils"
	"requests_responses"
)

func TestStoreLoad(t *testing.T) {
	var userDB myplaceutils.UserDB = make(myplaceutils.UserDB)
	var roomDB myplaceutils.RoomDB = make(myplaceutils.RoomDB)

	initUsers(userDB, 100, 10)
	initRooms(roomDB, 100, 100, 10, 10)
	var err error

	if err = StoreDBs(userDB, roomDB); err != nil {
		panic(err)
	}

	var userDB2 myplaceutils.UserDB
	var roomDB2 myplaceutils.RoomDB

	if userDB2, roomDB2, err = LoadDBs(); err != nil {
		panic(err)
	}

	// if equalRoomDBs.DeepEqual(roomDB, roomDB2) != true {
	// 	t.Errorf("roomDB != roomDB2")
	// }

	// if equalUserDBs.DeepEqual(userDB, userDB2) != true {
	// 	t.Errorf("userDB != userDB2")
	// }

	if len(roomDB) != len(roomDB2) {
		t.Errorf("len(roomDB) != len(roomDB2)")
	}

	if len(userDB) != len(userDB2) {
		t.Errorf("len(userDB) != len(userDB2)")
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

		users[u.UName] = u
	}
}

func initRooms(rooms myplaceutils.RoomDB, no_rooms, no_users int, no_msgs int, no_chans int) {
	for i := 0; i < no_rooms; i++ {
		var r myplaceutils.Room
		r.ID = i
		r.Name = fmt.Sprintf("room%v", i)

		r.Users = list.New()
		for j := 0; j < no_users; j++ {
			r.Users.PushBack(fmt.Sprintf("user%v", j))
		}

		r.Messages = make(map[int]myplaceutils.Message)
		for j := 0; j < no_msgs; j++ {
			msg := myplaceutils.Message{
				j,
				time.Now(),
				fmt.Sprintf("user%v", j),
				"msg body" }

			r.Messages[j] = msg
		}

		r.OutgoingChannels = list.New()
		for j := 0; j < no_chans; j++ {
			c := make(chan requests_responses.Response)
			r.OutgoingChannels.PushBack(c)
		}

		rooms[i] = r
	}
}

// func printRooms(rs []Room) {
// 	fmt.Println("\nRooms:");

// 	for _, r := range rs {
// 		fmt.Printf("\nName: %s\n", r.Name)
// 		fmt.Printf("No people: %v\n", r.NoPeople)
// 		fmt.Println("Messages:")

// 		for _, m := range r.Messages {
// 			fmt.Printf("Time: %v\n", m.Time)
// 			fmt.Printf("User: %s\n", m.Uname)
// 			fmt.Printf("Text: %s\n", m.Body)
// 		}
// 	}
// }

// func printUsers(us []User) {
// 	fmt.Println("\nUsers:")

// 	for _, u := range us {
// 		fmt.Printf("\nName: %s\n", u.Uname)
// 		fmt.Printf("Pass: %s\n", u.Pass)
// 		fmt.Printf("Rooms: %v\n", u.Rooms)
// 	}
// }
