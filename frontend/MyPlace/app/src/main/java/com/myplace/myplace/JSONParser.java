package com.myplace.myplace;


import android.support.annotation.NonNull;
import android.util.Log;
import android.widget.Toast;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.w3c.dom.ProcessingInstruction;

import java.util.ArrayList;
import java.util.Date;

import static com.myplace.myplace.models.RequestTypes.*;


/**
 * Created by alexis on 2017-04-28.
 */

public final class JSONParser {

    private static int number = 1;

    private static final String KEY_REQUEST_ID = "RequestID";
    private static final String KEY_USERNAME = "UName";
    private static final String KEY_PASSWORD = "Pass";
    private static final String KEY_SIGN_ACCEPTANCE = "Result";
    private static final String KEY_ROOM_ID = "RoomID";
    private static final String KEY_ROOM_NAME = "RoomName";
    private static final String KEY_ROOM_ID_ACCEPTED = "RoomIDAccepted";
    private static final String KEY_ROOM_LIST = "Rooms";
    private static final String KEY_ROOM_JOINED = "JoinedRoom";
    private static final String KEY_USER_LIST = "Users";
    private static final String KEY_MSG_ID = "MsgID";
    private static final String KEY_MESSAGE = "Msg";
    private static final String KEY_MSG_BODY = "Body";
    private static final String KEY_MSG_TIME  = "Time";
    private static final String KEY_LATEST_MSG = "LatestMsg";
    private static final String KEY_LATEST_MSG_ID = "LatestReadMsgID";
    private static final String KEY_MSG_LIST = "Messages";
    private static final String KEY_ERROR_CAUSE = "ErrorCause";

    private static final String TWO_CHAR_FORMAT = "%02d";
    private static final int NO_ID_FOUND = 0;



    public static int determineJSONType(String rawString) {
        int result = Integer.parseInt(rawString.substring(0, 2));
        if (result == ERROR_TYPE) throwErrorResponse(rawString);

        return result;
    }

    public static int determineRoomID(String rawString) {
        try {
            JSONObject json = makeProperJsonObject(rawString);
            return json.getInt(KEY_ROOM_ID);
        } catch (JSONException ignore) {
            return NO_ID_FOUND;
        }
    }

    public static int determineRequestID(String rawString) {
        try {
            JSONObject json = makeProperJsonObject(rawString);
            return json.getInt(KEY_REQUEST_ID);
        } catch (JSONException ignore) {
            return NO_ID_FOUND;
        }
    }


    public static String signupRequest(String userName, String password) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, userName);
        json.put(KEY_PASSWORD, password);

        return String.format(TWO_CHAR_FORMAT, SIGN_UP) + json.toString();
    }

    public static Boolean signupResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getBoolean(KEY_SIGN_ACCEPTANCE);
    }

    public static String signinRequest(String userName, String password) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, userName);
        json.put(KEY_PASSWORD, password);

        return String.format(TWO_CHAR_FORMAT, SIGN_IN) + json.toString();
    }

    public static Boolean signinResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getBoolean(KEY_SIGN_ACCEPTANCE);
    }

    public static String getRoomRequest(String userName) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, userName);

        return String.format(TWO_CHAR_FORMAT, GET_ROOMS) + json.toString();
    }

    public static ArrayList<RoomInfo> getRoomResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        if (!json.isNull(KEY_ROOM_LIST)) {
            JSONArray jsonRooms = json.getJSONArray(KEY_ROOM_LIST);
            ArrayList<RoomInfo> rooms = new ArrayList<>(jsonRooms.length());
            for (int i = 0; i < jsonRooms.length(); i++) {
                JSONObject r = jsonRooms.getJSONObject(i);

                rooms.add(constructRoom(r));
            }
            return rooms;
        }
        return null;
    }

    public static String getRoomUsersRequest(int roomID) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_ROOM_ID, roomID);

        return String.format(TWO_CHAR_FORMAT, GET_USERS) + json.toString();
    }

    public static ArrayList<String> getRoomUsersResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        JSONArray jsonUsers = json.getJSONArray(KEY_USER_LIST);
        ArrayList<String> usersList = new ArrayList<>(jsonUsers.length());
        for (int i = 0; i < jsonUsers.length(); i++) {
            String u = jsonUsers.getString(i);
            usersList.add(u);

        }
        return usersList;
    }


    public static String getOlderMsgsRequest(int roomID, int msgID) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_ROOM_ID, roomID);
        json.put(KEY_MSG_ID, msgID);

        return String.format(TWO_CHAR_FORMAT, GET_OLDER) + json.toString();
    }


    public static ArrayList<Message> getOlderMsgsResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        JSONArray jsonMsgs = json.getJSONArray(KEY_MSG_LIST);
        ArrayList<Message> messages = new ArrayList<>(jsonMsgs.length());
        for (int i = 0; i < jsonMsgs.length(); i++) {
            JSONObject m = jsonMsgs.getJSONObject(i);

            Message _msg = constructMessage(m);
            messages.add(_msg);
        }
        return messages;
    }

    public static String getNewerMsgsRequest(int roomID, int msgID) throws JSONException {
        String s = getOlderMsgsRequest(roomID, msgID);
        return s.replaceFirst(String.format(TWO_CHAR_FORMAT, GET_OLDER), String.format(TWO_CHAR_FORMAT, GET_NEWER));
    }

    public static ArrayList<Message> getNewerMsgsResponse(String rawString) throws JSONException {
        return getOlderMsgsResponse(rawString);
    }

    public static String joinRoomRequest(int roomID, String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_ROOM_ID, roomID);
        json.put(KEY_USERNAME, username);

        return String.format(TWO_CHAR_FORMAT, JOIN_ROOM) + json.toString();
    }

    public static RoomInfo joinRoomResponse(String rawString) throws JSONException, NullPointerException {
        JSONObject json = makeProperJsonObject(rawString);
        if (!json.getBoolean(KEY_ROOM_ID_ACCEPTED)) {
            return null;
        }
        return constructRoom(json.getJSONObject(KEY_ROOM_JOINED));
    }

    public static String leaveRoomRequest(int roomID, String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_ROOM_ID, roomID);
        json.put(KEY_USERNAME, username);

        return String.format(TWO_CHAR_FORMAT, LEAVE_ROOM) + json.toString();
    }

    public static int leaveRoomResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getInt(KEY_REQUEST_ID);
    }

    public static String createRoomRequest(String roomName, String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_ROOM_NAME, roomName);
        json.put(KEY_USERNAME, username);

        return String.format(TWO_CHAR_FORMAT, CREATE_ROOM) + json.toString();
    }

    public static RoomInfo createRoomResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        int _id = json.getInt(KEY_ROOM_ID);
        String _name = json.getString(KEY_ROOM_NAME);
        Room room = new Room(_id, _name);
        return new RoomInfo(room);
    }

    public static String postMsgRequest(Message msg) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, msg.getName());

        if (msg.getRoomID() == -1) throw new AssertionError();
        json.put(KEY_ROOM_ID, msg.getRoomID());
        json.put(KEY_MSG_BODY, msg.getText());

        return String.format(TWO_CHAR_FORMAT, MESSAGE) + json.toString();
    }

    public static Message messageRecieved(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return constructMessage(json.getJSONObject(KEY_MESSAGE));
    }

    public static String messageReadRequest(String username, int roomID, int msgID) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, username);
        json.put(KEY_ROOM_ID, roomID);
        json.put(KEY_MSG_ID, msgID);

        return String.format(TWO_CHAR_FORMAT, MSG_READ) + json.toString();
    }

    public static String signoutRequest(String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, username);

        return String.format(TWO_CHAR_FORMAT, SIGN_OUT) + json.toString();
    }

    public static String deleteUserRequest(String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put(KEY_USERNAME, username);

        return String.format(TWO_CHAR_FORMAT, DELETE_USER) + json.toString();
    }



    public static void throwErrorResponse(String rawString) {
        try {
            JSONObject json = makeProperJsonObject(rawString);
            String error = json.getString(KEY_ERROR_CAUSE);
            Log.e("JSONParser", "ERROR-RESPONSE: " + error);
            //throw new RuntimeException(error);
        } catch (JSONException e) {
            e.printStackTrace();
        }
    }







    @NonNull
    private static JSONObject makeProperJsonObject(String rawString) throws JSONException {
        return new JSONObject(rawString.substring(2));
    }

    @NonNull
    private static JSONObject constructJSONRequest() throws JSONException {
        JSONObject json = new JSONObject();

        json.put("RequestID", number++);
        return json;
    }

    private static Message constructMessage(JSONObject json) throws JSONException {
        int id = json.getInt(KEY_MSG_ID);
        int roomID = json.getInt(KEY_ROOM_ID);
        String fromUser = json.getString(KEY_USERNAME);
        long timestamp = json.getInt(KEY_MSG_TIME);
        String body = json.getString(KEY_MSG_BODY);

        return new Message(id, roomID, fromUser, body, timestamp);
    }

    private static RoomInfo constructRoom(JSONObject json) throws JSONException {
        int roomID = json.getInt(KEY_ROOM_ID);
        String roomName = json.getString(KEY_ROOM_NAME);
        Room room = new Room(roomID, roomName);

        JSONObject ms = json.getJSONObject(KEY_LATEST_MSG);
        Message msg = constructMessage(ms);

        int latestReadMsg = json.getInt(KEY_LATEST_MSG_ID);

        return new RoomInfo(room, msg, latestReadMsg);
    }
}
