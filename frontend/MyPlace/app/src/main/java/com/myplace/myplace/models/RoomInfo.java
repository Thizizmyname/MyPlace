package com.myplace.myplace.models;

/**
 * Created by alexis on 2017-05-09.
 */

public class RoomInfo {
    public Room room;
    public Message latestMsg;
    public int latestMsgRead;


    public RoomInfo(Room _room, Message msg, int latestRead) {
        room = _room;
        latestMsg = msg;
        latestMsgRead = latestRead;
    }

}
