package com.myplace.myplace;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-04-18.
 */

public class Room {
    private String roomName;
    public ArrayList<String> messagelist;

    public Room(String name) {
        roomName = name;
    }

    public String getName() {
        return roomName;
    }

    public String getLastMessage() {
        return "this is a message";
    }
}
