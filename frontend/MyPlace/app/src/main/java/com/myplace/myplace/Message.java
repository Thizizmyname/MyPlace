package com.myplace.myplace;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;

class Message {
    public String name = "Anders";
    String text = "";
    String date;
    private DateFormat df = new SimpleDateFormat("yy-MM-dd HH:mm");

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
