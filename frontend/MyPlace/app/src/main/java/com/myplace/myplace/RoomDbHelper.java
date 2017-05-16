package com.myplace.myplace;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;

import java.text.ParseException;
import java.util.ArrayList;
import java.util.Date;
import java.util.concurrent.ExecutionException;

/**
 * Created by jesper on 2017-04-26.
 */

public class RoomDbHelper extends SQLiteOpenHelper {
    private static final String DATABASE_NAME = "roomdb";
    private static final int DATABASE_VERSION = 1;

    private Context context;

    private static final String TABLE_ROOMS = "roomlist";


    public RoomDbHelper(Context context) {
        super(context, DATABASE_NAME,null, DATABASE_VERSION);
        this.context = context;
        SQLiteDatabase db = this.getWritableDatabase();
        onCreate(db);
        db.close();
    }

    public void onCreate(SQLiteDatabase db) {
        //Creates a table containing a list of all rooms and the last message in every room
        String CREATE_ROOMLIST_TABLE = "CREATE TABLE IF NOT EXISTS " + TABLE_ROOMS + "(roomname TEXT, lastmessage TEXT);";
        db.execSQL(CREATE_ROOMLIST_TABLE);
    }

    void createRoomTable(String roomName) {
        roomName = roomName.replace(" ", "_");

        SQLiteDatabase db = this.getWritableDatabase();
        final String CREATE_ROOM_TABLE = "CREATE TABLE IF NOT EXISTS "+roomName+"(name TEXT, message TEXT, date DATETIME);";
        db.execSQL(CREATE_ROOM_TABLE);
        db.close();
        if (!roomExists(roomName)) {
            addRoom(roomName);
        }
    }

    private boolean roomExists(String roomName) {
        roomName = roomName.replace(" ", "_");
        String query = "SELECT roomname FROM roomlist";
        SQLiteDatabase db = getReadableDatabase();

        Cursor c = db.rawQuery(query, null);
        if (c.moveToFirst()) {
            do {
                String selectedRoomName = c.getString(c.getColumnIndex("roomname"));
                if (selectedRoomName.equals(roomName)) {
                    return true;
                }
            } while (c.moveToNext());
            c.close();
        }
        return false;
    }

    private void addRoom(String roomName) {
        roomName = roomName.replace(" ", "_");

        ContentValues insertValues = new ContentValues();
        insertValues.put("roomname", roomName);
        insertValues.put("lastmessage", getLastMessage(roomName));

        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(TABLE_ROOMS, null, insertValues);
        db.close();

    }

    void deleteRoom(String roomName) {
        roomName = roomName.replace(" ", "_");

        SQLiteDatabase db = this.getWritableDatabase();
        db.delete(TABLE_ROOMS, "roomname = ?", new String[]{roomName});
        db.close();

        dropTable(roomName);
    }

    private void dropTable(String roomName) {
        roomName = roomName.replace(" ", "_");
        String query = "DROP TABLE IF EXISTS " + roomName;

        SQLiteDatabase db = this.getWritableDatabase();
        db.execSQL(query);
        db.close();
    }

    public void addMessage(String roomName, Message message) {
        roomName = roomName.replace(" ", "_");

        ContentValues insertValues = new ContentValues();
        insertValues.put("name", message.getName());
        insertValues.put("message", message.getText());
        insertValues.put("date", Message.df.format(message.getDate()));

        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(roomName, null, insertValues);
        db.close();
    }

    public String getLastMessage(String roomName) {
        try {
            roomName = roomName.replace(" ", "_");

            String query = "SELECT message FROM "+roomName;
            SQLiteDatabase db = getWritableDatabase();

            Cursor c = db.rawQuery(query, null);
            c.moveToLast();
            String message = c.getString(c.getColumnIndex("message"));

            c.close();
            db.close();
            return message;
        } catch (Exception ignore) {
            return context.getString(R.string.room_empty);
        }
    }

    public String getLastSender(String roomName) {
        try {
            roomName = roomName.replace(" ", "_");
            String query = "SELECT name FROM "+roomName;
            SQLiteDatabase db = getReadableDatabase();

            Cursor c = db.rawQuery(query, null);
            c.moveToLast();
            String name = c.getString(c.getColumnIndex("name"));

            c.close();
            db.close();
            return name;
        } catch (Exception ignore) {
            return "";
        }


    }

    ArrayList<Room> getRoomList(){
        ArrayList<Room> list=new ArrayList<>();
        String selectQuery = "SELECT * FROM "+TABLE_ROOMS;

        SQLiteDatabase db = this.getReadableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                String roomName = c.getString(c.getColumnIndex("roomname"));

                // TODO Replace ID to proper information
                Room room = new Room(0, roomName.replace("_", " "));
                list.add(room);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        return list;
    }

    ArrayList<Message> getMessages(String roomName){
        roomName = roomName.replace(" ", "_");

        ArrayList<Message> list=new ArrayList<>();
        String selectQuery = "SELECT * FROM "+roomName;

        SQLiteDatabase db = this.getReadableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                String name = c.getString(c.getColumnIndex("name"));
                String message = c.getString(c.getColumnIndex("message"));
                Date date = null;
                try {
                    date = Message.df.parse(c.getString(c.getColumnIndex("date")));
                } catch (ParseException e) {
                    e.printStackTrace();
                }

                Message newMessage = new Message(name, message, date);
                list.add(newMessage);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        return list;
    }

    private void dropRoomTable(String roomName) {
        roomName = roomName.replace(" ", "_");

        String query = "DROP TABLE IF EXISTS "+roomName;
        SQLiteDatabase db = this.getWritableDatabase();

        db.execSQL(query);
        db.close();
    }

    private void dropRoomListTable() {

        String query = "DROP TABLE IF EXISTS "+TABLE_ROOMS;
        SQLiteDatabase db = this.getWritableDatabase();

        db.execSQL(query);
        db.close();
    }

    public void dropAllTables() {
        String selectQuery = "SELECT * FROM "+TABLE_ROOMS;

        SQLiteDatabase db = this.getWritableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                String roomName = c.getString(c.getColumnIndex("roomname"));

                dropRoomTable(roomName);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        dropRoomListTable();
    }

    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {

    }
}
