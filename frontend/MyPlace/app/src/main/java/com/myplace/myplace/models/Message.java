package com.myplace.myplace.models;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;

public class Message {
    private int id;
    private int roomID;
    private String name;
    private String text;
    private long timestamp;
    public static DateFormat df = new SimpleDateFormat("HH:mm", Locale.getDefault());

    private static final int CROP_LIMIT = 32;

    public Message(int id, int roomID, String name, String text, long timestamp) {
        this.id = id;
        this.roomID = roomID;
        this.name = name;
        this.text = text;
        this.timestamp = timestamp;
    }

    public Message(int roomID, String name, String text, long timestamp) {
        this.roomID = roomID;
        this.name = name;
        this.text = text;
        this.timestamp = timestamp;
    }

/*    public Message(String name, String message, long timestamp) {
        this.name = name;
        this.text = message;
        this.timestamp = timestamp;
    }

    public Message(String username, String message) {
        this.name = username;
        this.text = message;
        this.timestamp = System.currentTimeMillis();

    }*/

    public int getId() {
        return id;
    }

    public int getRoomID() {
        return roomID;
    }

    public String getName() {
        return name;
    }

    public String getText() {
        return text;
    }

    public String getShortTime() {
        Date date = new Date(timestamp);
        return df.format(date);
    }

    public long getTimestamp() {
        return timestamp;
    }

    public String abbreviateText() {
        if (text.length() > CROP_LIMIT) {
            return text.substring(0, CROP_LIMIT-3) + "...";
        }
        else return text;
    }


}
