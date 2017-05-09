package com.myplace.myplace.models;

import org.json.JSONException;
import org.json.JSONObject;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;

public class Message {
    public int id;
    public int roomID;
    public String name;
    public String text;
    public Date date;
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

    public Message(String name, String message) {
        this.name = name;
        this.text = message;
        this.date = new Date();

    }

    public String getName() {
        return name;
    }

    public String abbreviateText() {
        if (text.length() > CROP_LIMIT) {
            return text.substring(0, CROP_LIMIT-3) + "...";
        }
        else return text;
    }
}
