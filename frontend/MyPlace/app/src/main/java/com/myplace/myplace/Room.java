package com.myplace.myplace;

import android.database.Cursor;
import android.os.Parcel;
import android.os.Parcelable;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-04-18.
 */

public class Room implements Parcelable {
    private String roomName;
    public ArrayList<Message> messageList = new ArrayList<>();

    public Room(String name) {
        roomName = name;
    }

    protected Room(Parcel in) {
        roomName = in.readString();
    }

    public static final Creator<Room> CREATOR = new Creator<Room>() {
        @Override
        public Room createFromParcel(Parcel in) {
            return new Room(in);
        }

        @Override
        public Room[] newArray(int size) {
            return new Room[size];
        }
    };

    public String getName() {
        return roomName;
    }

    public String getLastMessage() {
        //return "Detta är ett meddelande";
        try {
            String query = "SELECT message FROM "+roomName;
            Cursor c = MainActivity.roomDB.rawQuery(query, null);
            c.moveToLast();
            String message = c.getString(c.getColumnIndex("message"));
            c.close();
            return message;
            //return messageList.get(messageList.size()-1).text;
        } catch (Exception ignore) {
            return "This board has no messages";
        }

    }

    public void addMessage(Message msg) {
        messageList.add(msg);
    }








    @Override
    public int describeContents() {
        return 0;
    }

    @Override
    public void writeToParcel(Parcel dest, int flags) {
        dest.writeString(roomName);
    }
}
