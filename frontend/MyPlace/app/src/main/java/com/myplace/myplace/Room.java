package com.myplace.myplace;

import android.database.Cursor;
import android.os.Parcel;
import android.os.Parcelable;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-04-18.
 */

public class Room {
    private String roomName;
    public ArrayList<Message> messageList = new ArrayList<>();

    public Room(String name) {
        roomName = name;
    }


    public String getName() {
        return roomName;
    }

    public String getLastMessage() {
        return MainActivity.roomDB.getLastMessage(this.roomName);
    }

    public void addMessage(Message msg) {
        messageList.add(msg);
    }
}
