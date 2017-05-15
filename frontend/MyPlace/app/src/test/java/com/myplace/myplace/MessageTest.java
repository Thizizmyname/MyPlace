package com.myplace.myplace;

import com.myplace.myplace.models.Message;

import org.junit.Test;

import java.util.Date;

import static junit.framework.Assert.assertEquals;

/**
 * Created by jsepr on 2017-05-15.
 */

public class MessageTest {
    private static String TEST_NAME         = "Anders";
    private static String TEST_MESSAGE      = "Hej123";
    private static String TEST_LONG_MESSAGE = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas porta.";
    private static int    TEST_ID           = 1;
    private static int    TEST_ROOM_ID      = 123;


    @Test
    public void getIdTest() {
        Date date = new Date();
        Message message = new Message(TEST_ID, TEST_ROOM_ID, TEST_NAME, TEST_MESSAGE, date);

        assertEquals(message.getId(), TEST_ID);
    }

    @Test
    public void getRoomIdTest() {
        Date date = new Date();
        Message message = new Message(TEST_ID, TEST_ROOM_ID, TEST_NAME, TEST_MESSAGE, date);

        assertEquals(message.getRoomID(), TEST_ROOM_ID);
    }

    @Test
    public void getDateTest() {
        Date date = new Date();
        Message message = new Message(TEST_NAME, TEST_MESSAGE, date);

        assertEquals(message.getDate(), date);
    }

    @Test
    public void getNameTest() {
        Message message = new Message(TEST_NAME, TEST_MESSAGE);

        assertEquals(message.getName(), TEST_NAME);
    }

    @Test
    public void getTextTest() {
        Message message = new Message(TEST_NAME, TEST_MESSAGE);

        assertEquals(message.getText(), TEST_MESSAGE);
    }

    @Test
    public void abbreviateTextTest() {
        Message shortMessage = new Message(TEST_NAME, TEST_MESSAGE);
        Message longMessage = new Message(TEST_NAME, TEST_LONG_MESSAGE);

        assertEquals(shortMessage.abbreviateText().length(), TEST_MESSAGE.length());
        assertEquals(longMessage.abbreviateText().length(), 32);
    }
}
