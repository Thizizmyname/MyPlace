package com.myplace.myplace.models;

import android.content.Context;
import android.content.ContextWrapper;

import com.myplace.myplace.MainActivity;
import com.myplace.myplace.RoomDbHelper;
import com.myplace.myplace.models.Message;

/**
 * Created by alexis on 2017-04-18.
 */

public class Room {
    private int roomID;
    private String roomName;
    private Message latestMsg;




    public Room(String name) {
        roomName = name;
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
