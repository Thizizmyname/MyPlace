package com.myplace.myplace.models;

import android.content.Context;

/**
 * Created by alexis on 2017-05-09.
 */

public class RoomInfo {
    private Room room;
    private Message lastMsg = null;
    private int lastMsgRead;

    public final static String EMPTY_STRING = "";
    public final static String EMPTY_MESSAGE_STRING = "This room has no messages.";
    public final static long EMPTY_TIMESTAMP = 0;

    public RoomInfo(Room _room, Message msg, int latestRead) {
        room = _room;
        lastMsg = msg;
        lastMsgRead = latestRead;
    }

    public RoomInfo(Room _room) {
        room = _room;
        lastMsg = null;
        lastMsgRead = 0;
    }

    public int getLastMsgRead() {
        if (lastMsg != null) {
            return lastMsgRead;
        }
        return -1;
    }

    public Room getRoom() {
        return room;
    }

    public String getName() {
        return room.getName();
    }

    public int getRoomID() {
        return room.getRoomID();
    }

    public Message getLastMessage() {
        if (lastMsg != null) {
            return lastMsg;
        }
        return null;
    }

    public String getLastSender() {
        if (lastMsg != null) {
            return lastMsg.getName();
        }
        return EMPTY_STRING;
    }

    public String getLastMessageText() {
        if (lastMsg != null) {
            return lastMsg.getText();
        }
        return EMPTY_MESSAGE_STRING;
    }

    public long getLastMessageTime() {
        if (lastMsg != null) {
            return lastMsg.getTimestamp();
        }
        return EMPTY_TIMESTAMP;
    }

    public boolean hasLastMessage() {
        return lastMsg != null && lastMsg.getId() != -1;
    }
}
