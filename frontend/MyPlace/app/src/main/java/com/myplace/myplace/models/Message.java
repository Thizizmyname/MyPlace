package com.myplace.myplace.models;

import org.json.JSONException;
import org.json.JSONObject;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;

public class Message {
    private int id;
    private int roomID;
    private String name;
    private String text;
    private Date date;
    public static DateFormat df = new SimpleDateFormat("HH:mm", Locale.getDefault());

    private static final int CROP_LIMIT = 32;

    public Message(int id, int roomID, String name, String text, Date date) {
        this.id = id;
        this.roomID = roomID;
        this.name = name;
        this.text = text;
        this.date = date;
    }

    public Message(String name, String message, Date date) {
        this.name = name;
        this.text = message;
        this.date = date;
    }

    public Message(String username, String message) {
        this.name = username;
        this.text = message;
        this.date = new Date();

    }

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

    public Date getDate() {
        return date;
    }

    public String abbreviateText() {
        if (text.length() > CROP_LIMIT) {
            return text.substring(0, CROP_LIMIT-3) + "...";
        }
        else return text;
    }


}
