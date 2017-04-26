package com.myplace.myplace;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

import java.util.ArrayList;

/**
 * Created by jesper on 2017-04-26.
 */

class RoomDbHelper extends SQLiteOpenHelper {
    private static final String DATABASE_NAME = "Rooms.db";
    private static final int DATABASE_VERSION = 1;

    private static final String TABLE_ROOMS = "roomlist";


    RoomDbHelper(Context context) {
        super(context, DATABASE_NAME,null, DATABASE_VERSION);
    }

    public void onCreate(SQLiteDatabase db) {
        String CREATE_ROOMLIST_TABLE = "CREATE TABLE IF NOT EXISTS " + TABLE_ROOMS + "(roomname TEXT, lastmessage TEXT);";
        db.execSQL(CREATE_ROOMLIST_TABLE);
    }

    public void createRoomTable(String roomName) {
        SQLiteDatabase db = this.getWritableDatabase();
        final String CREATE_ROOM_TABLE = "CREATE TABLE IF NOT EXISTS "+roomName+"(name TEXT, message TEXT, date DATETIME);";
        db.execSQL(CREATE_ROOM_TABLE);
        db.close();
    }

    public void addRoom(String roomName) {

    }

    public void addMessage(String roomName, Message message) {
        ContentValues insertValues = new ContentValues();
        insertValues.put("name", message.getName());
        insertValues.put("message", message.text);
        insertValues.put("date", message.date);

        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(roomName, null, insertValues);
        db.close();
    }

    public String getLastMessage(String roomName) {
        String query = "SELECT message FROM "+roomName;
        SQLiteDatabase db = getWritableDatabase();

        Cursor c = db.rawQuery(query, null);
        c.moveToLast();
        String message = c.getString(c.getColumnIndex("message"));

        c.close();
        db.close();
        return message;
    }

    public ArrayList<Message> getMessages(String roomName){
        ArrayList<Message> list=new ArrayList<>();
        String selectQuery = "SELECT * FROM "+roomName;

        SQLiteDatabase db = this.getWritableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                String name = c.getString(c.getColumnIndex("name"));
                String message = c.getString(c.getColumnIndex("message"));
                String date = c.getString(c.getColumnIndex("date"));

                Message newMessage = new Message(name, message, date);
                list.add(newMessage);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        return list;
    }

    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {

    }
}
