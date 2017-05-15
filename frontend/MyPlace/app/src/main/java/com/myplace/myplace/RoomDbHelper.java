package com.myplace.myplace;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.CursorIndexOutOfBoundsException;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;

import java.text.ParseException;
import java.util.ArrayList;
import java.util.Date;
import java.util.concurrent.ExecutionException;

/**
 * Created by jesper on 2017-04-26.
 */

public class RoomDbHelper extends SQLiteOpenHelper {
    public static String DATABASE_NAME = "roomdb";
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
        String CREATE_ROOMLIST_TABLE = "CREATE TABLE IF NOT EXISTS " + TABLE_ROOMS + "(roomid INTEGER, roomname TEXT, lastmessage TEXT, lastmessagedate INTEGER);";
        db.execSQL(CREATE_ROOMLIST_TABLE);
    }

    void createRoomTable(int roomID, String roomName) {
        SQLiteDatabase db = this.getWritableDatabase();
        String roomIDString = "_" + roomID;
        final String CREATE_ROOM_TABLE = "CREATE TABLE IF NOT EXISTS "+roomIDString+"(messageid INTEGER, name TEXT, text TEXT, date INTEGER);";
        db.execSQL(CREATE_ROOM_TABLE);
        db.close();
        if (!roomExists(roomID)) {
            addRoom(roomID, roomName);
        }
    }

    private boolean roomExists(int roomID) {
        String query = "SELECT roomid FROM roomlist";
        SQLiteDatabase db = getReadableDatabase();

        Cursor c = db.rawQuery(query, null);
        if (c.moveToFirst()) {
            do {
                int selectedID = c.getInt(c.getColumnIndex("roomid"));
                if (selectedID == roomID) {
                    return true;
                }
            } while (c.moveToNext());
            c.close();
        }
        return false;
    }

    private void addRoom(int roomID, String roomName) {
        ContentValues insertValues = new ContentValues();
        insertValues.put("roomid", roomID);
        insertValues.put("roomname", roomName);
        Message lastmessage;
        try {
            lastmessage = getLastMessage(roomID);
            insertValues.put("lastmessage", lastmessage.getText());
            insertValues.put("lastmessagedate", lastmessage.getDate().getTime());
        } catch (Exception e) {
            insertValues.put("lastmessage", "");
            insertValues.put("lastmessagedate", 0);
        }

        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(TABLE_ROOMS, null, insertValues);
        db.close();
    }

    void deleteRoom(int roomID) {
        SQLiteDatabase db = this.getWritableDatabase();
        String roomIDString = "_" + roomID;
        db.delete(TABLE_ROOMS, "roomid = ?", new String[]{roomIDString});
        db.close();

        dropTable(roomID);
    }

    private void dropTable(int roomID) {
        String roomIDString = "_"+roomID;
        String query = "DROP TABLE IF EXISTS "+roomIDString;

        SQLiteDatabase db = this.getWritableDatabase();
        db.execSQL(query);
        db.close();
    }

    void addMessage(int roomID, Message message) {
        ContentValues insertValues = new ContentValues();
        insertValues.put("messageid", message.getId());
        insertValues.put("name", message.getName());
        insertValues.put("text", message.getText());
        insertValues.put("date", Message.df.format(message.getDate()));

        String roomIDString = "_" + roomID;
        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(roomIDString, null, insertValues);
        db.close();
    }

    public Message getLastMessage(int roomID) throws Exception {
        String roomIDString = "_" + roomID;
        String query = "SELECT * FROM "+roomIDString;
        SQLiteDatabase db = getWritableDatabase();

        Cursor c = db.rawQuery(query, null);
        c.moveToLast();
        int messageID = c.getInt(c.getColumnIndex("messageid"));
        String name = c.getString(c.getColumnIndex("name"));
        String text = c.getString(c.getColumnIndex("text"));
        int date = c.getInt(c.getColumnIndex("date"));

        Message message = new Message(messageID, roomID, name, text, new Date(date));
        c.close();
        db.close();
        return message;
    }

//    public String getLastSender(String roomID) {
//        try {
//            String query = "SELECT name FROM ["+roomID+"]";
//            SQLiteDatabase db = getReadableDatabase();
//
//            Cursor c = db.rawQuery(query, null);
//            c.moveToLast();
//            String name = c.getString(c.getColumnIndex("name"));
//
//            c.close();
//            db.close();
//            return name;
//        } catch (Exception ignore) {
//            return "";
//        }
//
//
//    }

    ArrayList<RoomInfo> getRoomList(){
        ArrayList<RoomInfo> list=new ArrayList<>();
        String selectQuery = "SELECT * FROM "+TABLE_ROOMS;

        SQLiteDatabase db = this.getReadableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                int roomID      = c.getInt(c.getColumnIndex("roomid"));
                String roomName = c.getString(c.getColumnIndex("roomname"));

                // TODO Replace ID to proper information
                Room room = new Room(roomID, roomName);

                Message lastMessage;
                try {
                    lastMessage = getLastMessage(roomID);
                } catch (Exception e) {
                    lastMessage = null;
                }

                RoomInfo rInfo = new RoomInfo(room, lastMessage, 0);
                list.add(rInfo);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        return list;
    }

    ArrayList<Message> getMessages(int roomID){
        ArrayList<Message> list=new ArrayList<>();
        String roomIDString = "_" + roomID;
        String selectQuery = "SELECT * FROM "+roomIDString;

        SQLiteDatabase db = this.getReadableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                int messageID = c.getInt(c.getColumnIndex("messageid"));
                String name = c.getString(c.getColumnIndex("name"));
                String message = c.getString(c.getColumnIndex("text"));
                Date date = null;
                try {
                    date = new Date(c.getInt(c.getColumnIndex("date")));
                } catch (NullPointerException e) {
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

    private void dropRoomTable(int roomID) {
        String roomIDString = "_" + roomID;
        String query = "DROP TABLE IF EXISTS "+roomIDString;
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
                int roomID = c.getInt(c.getColumnIndex("roomid"));

                dropRoomTable(roomID);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        dropRoomListTable();
    }

    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {

    }
}
