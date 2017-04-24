package com.myplace.myplace;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;

class Message {
    public String name = "Anders";
    public String text = "";
    public String date;
    private DateFormat df = new SimpleDateFormat("HH:mm");

    Message(String name, String message, String date) {
        this.name = name;
        this.text = message;
        this.date = date;
    }

    Message(String name, String message) {
        this.name = name;
        this.text = message;
        Date mDate = new Date();
        this.date = df.format(mDate);
    }

    public String getName() {
        return name;
    }
}
