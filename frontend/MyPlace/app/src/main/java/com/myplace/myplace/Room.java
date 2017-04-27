package com.myplace.myplace;

import android.database.Cursor;
import android.os.Parcel;
import android.os.Parcelable;

import java.util.ArrayList;

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

    public String getLastMessage() {
        return MainActivity.roomDB.getLastMessage(this.roomName);
    }

}
