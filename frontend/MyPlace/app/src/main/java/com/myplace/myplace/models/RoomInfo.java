package com.myplace.myplace.models;

import android.content.Context;

/**
 * Created by alexis on 2017-05-09.
 */

public class RoomInfo {
    private Room room;
    private Message latestMsg;
    private int latestMsgRead;

    public RoomInfo(Room _room, Message msg, int latestRead) {
        room = _room;
        latestMsg = msg;
        latestMsgRead = latestRead;
    }

    public RoomInfo(Room _room) {
        room = _room;
        latestMsg = null;
        latestMsgRead = 0;
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

    public Message getLatestMessage(Context ctx) {
        Message message = null;
        if (latestMsg != null) {
            message = latestMsg;
        } else {
            try {
                message = room.getLastMessage(ctx);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        return message;
    }

    public String getLastSender(Context ctx) {
        return room.getLastSender(ctx);
    }

    public String getLastMessageText(Context ctx) {
        return room.getLastMessageText(ctx);
    }
}
