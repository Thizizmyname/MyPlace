package com.myplace.myplace;

/**
 * Created by jsepr on 2017-05-15.
 */

import android.content.Context;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;

import static android.support.test.InstrumentationRegistry.getContext;
import static android.support.test.InstrumentationRegistry.getInstrumentation;
import static android.support.test.InstrumentationRegistry.getTargetContext;
import static junit.framework.Assert.assertEquals;

import android.support.test.InstrumentationRegistry;
import android.support.test.runner.AndroidJUnit4;
import android.test.RenamingDelegatingContext;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;

/**
 * Created by jsepr on 2017-05-15.
 */
@RunWith(AndroidJUnit4.class)
public class RoomInstrumentedTest {
    private static final String TEST_DATABASE = "test_db";
    private static int    TEST_ROOM_ID   = 123;
    private static String TEST_ROOM_NAME = "TestRoom";
    private static String TEST_NAME         = "Anders";
    private static String TEST_MESSAGE      = "Hej123";

    private Context testContext;

    private RoomDbHelper db = null;

    @Before
    public void setup() {
        RoomDbHelper.DATABASE_NAME = TEST_DATABASE;
        Message testMessage = new Message(TEST_NAME, TEST_MESSAGE);

        testContext = getInstrumentation().getTargetContext();
        db = new RoomDbHelper(testContext);

        db.createRoomTable(TEST_ROOM_NAME);
        db.addMessage(TEST_ROOM_NAME, testMessage);
    }

    @Test
    public void getLastMessageTest() {
        Room room = new Room(TEST_ROOM_ID, TEST_ROOM_NAME);
        String lastMessage = room.getLastMessage(testContext);
        assertEquals(lastMessage, TEST_MESSAGE);
    }

    @Test
    public void getLastSenderTest() {
        Room room = new Room(TEST_ROOM_ID, TEST_ROOM_NAME);
        String lastSender = room.getLastSender(testContext);
        assertEquals(lastSender, TEST_NAME);
    }

    @After
    public void teardown() {
        RoomDbHelper.DATABASE_NAME = TEST_DATABASE;
        db.dropAllTables();
    }
}