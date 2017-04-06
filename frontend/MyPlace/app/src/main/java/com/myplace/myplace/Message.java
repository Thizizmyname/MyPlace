package com.myplace.myplace;

import java.util.Date;

/**
 * Created by jesper on 2017-04-06.
 */

class Message {
    String name;
    String text;
    String date;

    Message(String name, String message) {
        this.name = name;
        this.text = message;
        this.date = new Date().toString();
    }
}
