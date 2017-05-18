package com.myplace.myplace.models;

import android.content.Context;
import android.database.CursorIndexOutOfBoundsException;

import com.myplace.myplace.R;
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

    public int getRoomID() {
        return roomID;
    }

    public String getName() {
        return roomName;
    }

    public String getLastSender(Context ctx) {
        try {
            return getLastMessage(ctx).getName();
        } catch (Exception e) {
            return "";
        }
    }

    public String getLastMessageText(Context ctx) {
        try {
            return getLastMessage(ctx).getText();
        } catch (Exception e) {
            return ctx.getResources().getString(R.string.room_empty);
        }
    }

    public Message getLastMessage(Context ctx) throws Exception {
        RoomDbHelper roomDB = new RoomDbHelper(ctx);
        return roomDB.getLastMessage(this.roomID);
    }


}
