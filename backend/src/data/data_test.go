package data

import (
	"fmt"
	"time"
	"math/rand"
	//"encoding/json"
	"reflect"
	"testing"
)

func TestStoreLoad(t *testing.T) {
	userDB := make(UserDB, 10)
	roomDB := make(RoomDB, 10)

	//math.Seed(time.Now().Unix())
	initUsers(userDB)
	initRooms(roomDB)
	var err error

	// printRooms(roomsDB)
	// printUsers(usersDB)

	// roomsDBJ, e1 := json.Marshal(roomsDB)
	// if e1 != nil {
	// 	panic(e1)
	// }

	// usersDBJ, e2 := json.Marshal(usersDB)
	// if e2 != nil {
	// 	panic(e2)
	// }

	//fmt.Println(string(roomsDBJ))

	if err = StoreDBs(userDB, roomDB); err != nil {
		panic(err)
	}

	// var roomsDB2 []Room
	// var usersDB2 []User

	// if e := json.Unmarshal(roomsDBJ, &roomsDB2); e != nil {
	// 	panic(e)
	// }

	// if e := json.Unmarshal(usersDBJ, &usersDB2); e != nil {
	// 	panic(e)
	// }

	var userDB2 UserDB
	var roomDB2 RoomDB

	if userDB2, roomDB2, err = LoadDBs(); err != nil {
		panic(err)
	}

	if reflect.DeepEqual(roomDB, roomDB2) != true {
		t.Errorf("roomDB != roomDB2")
	}

	if reflect.DeepEqual(userDB, userDB2) != true {
		t.Errorf("userDB != userDB2")
	}
}

func initUsers(users UserDB) {
	for i := 0; i < len(users); i++ {
		users[i].Uname = fmt.Sprintf("user%v", i)
		users[i].Pass = fmt.Sprintf("pass%v", i)
		rs := make([]string, rand.Intn(5))

		for j := 0; j < len(rs); j++ {
			rs[j] = fmt.Sprintf("room%v", j)
		}

		users[i].Rooms = rs
	}
}

func initRooms(rooms RoomDB) {
	for i := 0; i < len(rooms); i++ {
		rooms[i].Name = fmt.Sprintf("room%v", i)
		rooms[i].NoPeople = rand.Intn(50)
		msgs := make([]Message, rand.Intn(20))

		for j := 0; j < len(msgs); j++ {
			msgs[j].Time = time.Now()
			msgs[j].Uname = "some user"
			msgs[j].Text = fmt.Sprintf("message%v", j)
		}
		rooms[i].Messages = msgs
	}
}

func printRooms(rs []Room) {
	fmt.Println("\nRooms:");

	for _, r := range rs {
		fmt.Printf("\nName: %s\n", r.Name)
		fmt.Printf("No people: %v\n", r.NoPeople)
		fmt.Println("Messages:")

		for _, m := range r.Messages {
			fmt.Printf("Time: %v\n", m.Time)
			fmt.Printf("User: %s\n", m.Uname)
			fmt.Printf("Text: %s\n", m.Text)
		}
	}
}

func printUsers(us []User) {
	fmt.Println("\nUsers:")

	for _, u := range us {
		fmt.Printf("\nName: %s\n", u.Uname)
		fmt.Printf("Pass: %s\n", u.Pass)
		fmt.Printf("Rooms: %v\n", u.Rooms)
	}
}

// func TestReverse(t *testing.T) {
// 	cases := []struct {
// 		in, want string
// 	} {
// 		{"Hello, world", "dlrow ,olleH"},
// 		{"Hello, 世界", "界世 ,olleH"},
// 		{"", ""},
// 	}
// 	for _, c := range cases {
// 		got := Reverse(c.in)
// 		if got != c.want {
// 			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
// 		}
// 	}
// }
