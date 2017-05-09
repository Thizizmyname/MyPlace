package com.myplace.myplace.models;

import android.content.Context;

import com.myplace.myplace.RoomDbHelper;

import org.json.JSONException;
import org.json.JSONObject;

/**
 * Created by alexis on 2017-04-18.
 */

public class Room {
    private int roomID;
    private String roomName;





    public Room(int _roomID, String _roomName) {
        roomID = _roomID;
        roomName = _roomName;
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
