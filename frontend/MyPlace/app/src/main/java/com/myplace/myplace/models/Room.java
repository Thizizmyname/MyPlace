package com.myplace.myplace.models;

import android.content.Context;

import com.myplace.myplace.RoomDbHelper;

/**
 * Created by alexis on 2017-04-18.
 */

public class Room {
    private int roomID;
    private String roomName;
    private Message latestMsg;
    private int latestReadMsg;




    public Room(String name) {
        roomName = name;
    }

    public Room(int _roomID, String _roomName, Message _latestMsg, int _latestReadMsg) {
        roomID = _roomID;
        roomName = _roomName;
        latestMsg = _latestMsg;
        latestReadMsg = _latestReadMsg;
    }


    public String getName() {
        return roomName;
    }

    public String getLastMessage(Context ctx) {
        RoomDbHelper roomDB = new RoomDbHelper(ctx);
        return roomDB.getLastMessage(this.roomName);
        //return MainActivity.roomDB.getLastMessage(this.roomName);
    }

}
