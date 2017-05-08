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

    public Message(JSONObject json) throws JSONException {
        id = json.getInt("MsgID");
        roomID = json.getInt("RoomID");
        name = json.getString("FromUser");
        date = new Date(json.getInt("Time"));
        text = json.getString("Body");

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
}
