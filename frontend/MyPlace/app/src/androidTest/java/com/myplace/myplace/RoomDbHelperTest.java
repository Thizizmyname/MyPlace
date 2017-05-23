package com.myplace.myplace;

import android.content.Context;

import com.myplace.myplace.models.Message;

import org.junit.Before;

import static android.support.test.InstrumentationRegistry.getInstrumentation;

/**
 * Created by jsepr on 2017-05-15.
 */

public class RoomDbHelperTest {
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

        db.createRoomTable(TEST_ROOM_ID, TEST_ROOM_NAME);
        db.addMessage(TEST_ROOM_ID, testMessage);
    }
}
