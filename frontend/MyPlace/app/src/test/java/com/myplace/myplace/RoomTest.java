package com.myplace.myplace;

import com.myplace.myplace.models.Room;

import org.junit.Test;
import static org.junit.Assert.assertEquals;

/**
 * Created by jsepr on 2017-05-15.
 */

public class RoomTest {
    private static int    TEST_ROOM_ID   = 123;
    private static String TEST_ROOM_NAME = "TestRoom";

    @Test
    public void getNameTest() {
        Room room = new Room(TEST_ROOM_ID, TEST_ROOM_NAME);

        assertEquals(room.getName(), TEST_ROOM_NAME);
    }

    @Test
    public void getRoomIdTest() {
        Room room = new Room(TEST_ROOM_ID, TEST_ROOM_NAME);

        assertEquals(room.getRoomID(), TEST_ROOM_ID);
    }

}
