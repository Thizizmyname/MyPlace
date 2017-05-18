package com.myplace.myplace;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;
import android.util.Log;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;

import java.util.ArrayList;
import java.util.Date;

/**
 * Created by jesper on 2017-04-26.
 */

public class RoomDbHelper extends SQLiteOpenHelper {
    public static String DATABASE_NAME = "roomdb";
    private static final int DATABASE_VERSION = 1;

    // roomlist-table strings
    private static final String ROOMLIST_ROOMID = "roomid";
    private static final String ROOMLIST_ROOMNAME = "roomname";
    private static final String ROOMLIST_LASTMESSAGE = "lastmessage";
    private static final String ROOMLIST_LM_TIMESTAMP = "lmtimestamp";

    // roomtables strings
    private static final String ROOM_MESSAGEID = "messageid";
    private static final String ROOM_SENDER = "name";
    private static final String ROOM_TEXT = "text";
    private static final String ROOM_TIMESTAMP = "timestamp";

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
        String CREATE_ROOMLIST_TABLE = "CREATE TABLE IF NOT EXISTS " + TABLE_ROOMS + "("
                + ROOMLIST_ROOMID       +" INTEGER, "
                + ROOMLIST_ROOMNAME     +" TEXT, "
                + ROOMLIST_LASTMESSAGE  +" TEXT, "
                + ROOMLIST_LM_TIMESTAMP +" INTEGER);";

        db.execSQL(CREATE_ROOMLIST_TABLE);
    }

    private String getRoomIdString (int roomID) {
        return "_"+roomID;
    }

    void createRoomTable(int roomID, String roomName) {
        SQLiteDatabase db = this.getWritableDatabase();
        String roomIDString = getRoomIdString(roomID);

        final String CREATE_ROOM_TABLE = "CREATE TABLE IF NOT EXISTS "+roomIDString+"("
                +ROOM_MESSAGEID + " INTEGER, "
                +ROOM_SENDER    + " TEXT, "
                +ROOM_TEXT      + " TEXT, "
                +ROOM_TIMESTAMP + " INTEGER);";

        db.execSQL(CREATE_ROOM_TABLE);
        db.close();
        if (!roomExists(roomID)) {
            addRoom(roomID, roomName);
        }
    }

    private boolean roomExists(int roomID) {
        String query = "SELECT "+ROOMLIST_ROOMID+" FROM "+TABLE_ROOMS;
        SQLiteDatabase db = getReadableDatabase();

        Cursor c = db.rawQuery(query, null);
        if (c.moveToFirst()) {
            do {
                int selectedID = c.getInt(c.getColumnIndex(ROOMLIST_ROOMID));
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
        insertValues.put(ROOMLIST_ROOMID, roomID);
        insertValues.put(ROOMLIST_ROOMNAME, roomName);
        Message lastmessage;
        try {
            lastmessage = getLastMessage(roomID);
            insertValues.put(ROOMLIST_LASTMESSAGE, lastmessage.getText());
            insertValues.put(ROOMLIST_LM_TIMESTAMP, lastmessage.getTimestamp());
        } catch (Exception e) {
            insertValues.put(ROOMLIST_LASTMESSAGE, "");
            insertValues.put(ROOMLIST_LM_TIMESTAMP, 0);
        }

        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(TABLE_ROOMS, null, insertValues);
        db.close();
    }

    void deleteRoom(int roomID) {
        SQLiteDatabase db = this.getWritableDatabase();
        db.delete(TABLE_ROOMS, ROOMLIST_ROOMID+" = ?", new String[]{Integer.toString(roomID)});
        db.close();

        dropTable(roomID);
    }

    private void dropTable(int roomID) {
        String roomIDString = getRoomIdString(roomID);
        String query = "DROP TABLE IF EXISTS "+roomIDString;

        SQLiteDatabase db = this.getWritableDatabase();
        db.execSQL(query);
        db.close();
    }

    public void addMessage(int roomID, Message message) {
        ContentValues insertValues = new ContentValues();
        insertValues.put(ROOM_MESSAGEID, message.getId());
        insertValues.put(ROOM_SENDER, message.getName());
        insertValues.put(ROOM_TEXT, message.getText());
        long timestamp = message.getTimestamp();
        insertValues.put(ROOM_TIMESTAMP, timestamp);

        String roomIDString = getRoomIdString(roomID);
        SQLiteDatabase db = this.getWritableDatabase();
        db.insert(roomIDString, null, insertValues);

        ContentValues updatedLastMessage = new ContentValues();
        updatedLastMessage.put(ROOMLIST_LM_TIMESTAMP, timestamp);
        db.update(TABLE_ROOMS, updatedLastMessage, ROOMLIST_ROOMID+ "= ?", new String[]{Integer.toString(roomID)});
        db.close();
    }

    public Message getLastMessage(int roomID) throws Exception {
        String roomIDString = getRoomIdString(roomID);
        String query = "SELECT * FROM "+roomIDString;
        SQLiteDatabase db = getWritableDatabase();

        Cursor c = db.rawQuery(query, null);
        c.moveToLast();
        int messageID = c.getInt(c.getColumnIndex(ROOM_MESSAGEID));
        String name = c.getString(c.getColumnIndex(ROOM_SENDER));
        String text = c.getString(c.getColumnIndex(ROOM_TEXT));
        long timestamp = c.getLong(c.getColumnIndex(ROOM_TIMESTAMP));

        Message message = new Message(messageID, roomID, name, text, timestamp);
        c.close();
        db.close();
        return message;
    }

    ArrayList<RoomInfo> getRoomList(){
        ArrayList<RoomInfo> list=new ArrayList<>();
        String selectQuery = "SELECT * FROM "+TABLE_ROOMS+" ORDER BY "+ROOMLIST_LM_TIMESTAMP+" DESC";

        SQLiteDatabase db = this.getReadableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            int i = 0;
            do {
                int roomID      = c.getInt(c.getColumnIndex(ROOMLIST_ROOMID));
                String roomName = c.getString(c.getColumnIndex(ROOMLIST_ROOMNAME));

                Room room = new Room(roomID, roomName);

                Message lastMessage;
                try {
                    lastMessage = getLastMessage(roomID);
                } catch (Exception e) {
                    lastMessage = null;
                }
                Log.d("loop iteration", Integer.toString(i));
                RoomInfo rInfo = new RoomInfo(room, lastMessage, 0);
                list.add(rInfo);
                ++i;
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        return list;
    }

    ArrayList<Message> getMessages(int roomID){
        ArrayList<Message> list=new ArrayList<>();
        String roomIDString = getRoomIdString(roomID);
        String selectQuery = "SELECT * FROM "+roomIDString+" ORDER BY "+ ROOM_TIMESTAMP +" ASC";

        SQLiteDatabase db = this.getReadableDatabase();
        Cursor c = db.rawQuery(selectQuery,null);

        if (c.moveToFirst()) {
            do {
                int messageID = c.getInt(c.getColumnIndex(ROOM_MESSAGEID));
                String name = c.getString(c.getColumnIndex(ROOM_SENDER));
                String message = c.getString(c.getColumnIndex(ROOM_TEXT));
                long timestamp = c.getLong(c.getColumnIndex(ROOM_TIMESTAMP));

                Message newMessage = new Message(messageID, roomID, name, message, timestamp);
                list.add(newMessage);
            } while (c.moveToNext());
            c.close();
        }
        db.close();
        return list;
    }

    private void dropRoomTable(int roomID) {
        String roomIDString = getRoomIdString(roomID);
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
                int roomID = c.getInt(c.getColumnIndex(ROOMLIST_ROOMID));

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
