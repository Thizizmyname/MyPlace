package com.myplace.myplace;


import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-04-28.
 */

public class JSONParser {

    private static int number;

    public static String signupRequest(String userName, String password) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("UserName", userName);
        json.put("Password", password);

        return json.toString();
    }

    public static Boolean signupResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        return json.getBoolean("Result");
    }

    public static String signinRequest(String userName, String password) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("UserName", userName);
        json.put("Password", password);

        return json.toString();
    }

    public static Boolean signinResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        return json.getBoolean("Result");
    }

    public static String getRoomRequest(String userName) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("UserName", userName);

        return json.toString();
    }

    public static ArrayList<Room> getRoomResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        JSONArray jsonRooms = json.getJSONArray("Rooms");
        ArrayList<Room> rooms = new ArrayList<>(jsonRooms.length());
        for (int i = 0; i < jsonRooms.length(); i++) {
            JSONObject r = jsonRooms.getJSONObject(i);

            rooms.add(new Room(r));
        }
        return rooms;
    }


    public static String getRoomUsersRequest(int roomID) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("RoomID", roomID);

        return json.toString();
    }


    public static ArrayList<String> getRoomUsersResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        JSONArray jsonUsers = json.getJSONArray("Users");
        ArrayList<String> usersList = new ArrayList<>(jsonUsers.length());
        for (int i = 0; i < jsonUsers.length(); i++) {
            String u = jsonUsers.getString(i);
            usersList.add(u);

        }
        return usersList;
    }

    public static String getOlderMsgsRequest(int roomID, int msgID) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("RoomID", roomID);
        json.put("MsgID", msgID);

        return json.toString();
    }

    public static ArrayList<Message> getOlderMsgsResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        JSONArray jsonMsgs = json.getJSONArray("Messages");
        ArrayList<Message> messages = new ArrayList<>(jsonMsgs.length());
        for (int i = 0; i < jsonMsgs.length(); i++) {
            JSONObject m = jsonMsgs.getJSONObject(i);

            Message _msg = new Message(m);
            messages.add(_msg);
        }
        return messages;
    }

    public static String getNewerMsgsRequest(int roomID, int msgID) throws JSONException {
        return getOlderMsgsRequest(roomID, msgID);
    }

    public static ArrayList<Message> getNewerMsgsResponse(String rawString) throws JSONException {
        return getOlderMsgsResponse(rawString);
    }


    public static String joinRoomRequest(int roomID, String username) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("RoomID", roomID);
        json.put("Username", username);

        return json.toString();
    }

    public static Room joinRoomResponse(String rawString) throws JSONException, NullPointerException {
        JSONObject json = new JSONObject(rawString);
        if (!json.getBoolean("RoomIDAccepted")) {
            return null;
        }
        return new Room(json.getJSONObject("JoinedRoom"));
    }

    public static String leaveRoomRequest(int roomID, String username) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("RoomID", roomID);
        json.put("Username", username);

        return json.toString();
    }

    public static int leaveRoomResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        return json.getInt("RequestID");
    }

    public static String createRoomRequest(String roomName, String username) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("RoomName", roomName);
        json.put("Username", username);

        return json.toString();
    }

    public static int createRoomResponse(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        return json.getInt("RoomID");
    }

    public static String postMsgRequest(String username, Message msg) throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        json.put("Username", username);

        if (msg.roomID == 0) throw new AssertionError();
        json.put("RoomID", msg.roomID);
        json.put("RoomID", msg.text);

        return json.toString();
    }

    public static Message messageRecieved(String rawString) throws JSONException {
        JSONObject json = new JSONObject(rawString);
        return  new Message(json.getJSONObject("Msg"));
    }


}
