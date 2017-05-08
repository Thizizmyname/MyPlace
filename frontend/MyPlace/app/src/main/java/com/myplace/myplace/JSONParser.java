package com.myplace.myplace;


import android.support.annotation.NonNull;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

import static com.myplace.myplace.models.RequestTypes.*;


/**
 * Created by alexis on 2017-04-28.
 */

public final class JSONParser {

    private static int number;

    private static final String TWO_CHAR_FORMAT = "%03d";



    public static int determineJSONType(String rawString) {
        int result = Integer.parseInt(rawString.substring(0, 1));
        if (result == ERROR_TYPE) readErrorResponse(rawString);

        return result;
    }

    public static int determineRoomID(String rawString) {
        try {
            JSONObject json = makeProperJsonObject(rawString);
            return json.getInt("RoomID");
        } catch (JSONException ignore) {
            return 0;
        }
    }

    public static int determineRequestID(String rawString) {
        try {
            JSONObject json = makeProperJsonObject(rawString);
            return json.getInt("RequestID");
        } catch (JSONException ignore) {
            return 0;
        }
    }


    public static String signupRequest(String userName, String password) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("UserName", userName);
        json.put("Password", password);

        return String.format(TWO_CHAR_FORMAT, SIGN_UP) + json.toString();
    }

    public static Boolean signupResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getBoolean("Result");
    }

    public static String signinRequest(String userName, String password) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("UserName", userName);
        json.put("Password", password);

        return String.format(TWO_CHAR_FORMAT, SIGN_IN) + json.toString();
    }

    public static Boolean signinResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getBoolean("Result");
    }

    public static String getRoomRequest(String userName) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("UserName", userName);

        return String.format(TWO_CHAR_FORMAT, GET_ROOMS) + json.toString();
    }

    public static ArrayList<Room> getRoomResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        JSONArray jsonRooms = json.getJSONArray("Rooms");
        ArrayList<Room> rooms = new ArrayList<>(jsonRooms.length());
        for (int i = 0; i < jsonRooms.length(); i++) {
            JSONObject r = jsonRooms.getJSONObject(i);

            rooms.add(new Room(r));
        }
        return rooms;
    }

    public static String getRoomUsersRequest(int roomID) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("RoomID", roomID);

        return String.format(TWO_CHAR_FORMAT, GET_USERS) + json.toString();
    }

    public static ArrayList<String> getRoomUsersResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        JSONArray jsonUsers = json.getJSONArray("Users");
        ArrayList<String> usersList = new ArrayList<>(jsonUsers.length());
        for (int i = 0; i < jsonUsers.length(); i++) {
            String u = jsonUsers.getString(i);
            usersList.add(u);

        }
        return usersList;
    }


    public static String getOlderMsgsRequest(int roomID, int msgID) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("RoomID", roomID);
        json.put("MsgID", msgID);

        return String.format(TWO_CHAR_FORMAT, GET_OLDER) + json.toString();
    }


    public static ArrayList<Message> getOlderMsgsResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
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
        String s = getOlderMsgsRequest(roomID, msgID);
        return s.replaceFirst(String.format(TWO_CHAR_FORMAT, GET_OLDER), String.format(TWO_CHAR_FORMAT, GET_NEWER));
    }

    public static ArrayList<Message> getNewerMsgsResponse(String rawString) throws JSONException {
        return getOlderMsgsResponse(rawString);
    }

    public static String joinRoomRequest(int roomID, String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("RoomID", roomID);
        json.put("Username", username);

        return String.format(TWO_CHAR_FORMAT, JOIN_ROOM) + json.toString();
    }

    public static Room joinRoomResponse(String rawString) throws JSONException, NullPointerException {
        JSONObject json = makeProperJsonObject(rawString);
        if (!json.getBoolean("RoomIDAccepted")) {
            return null;
        }
        return new Room(json.getJSONObject("JoinedRoom"));
    }


    public static String leaveRoomRequest(int roomID, String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("RoomID", roomID);
        json.put("Username", username);

        return String.format(TWO_CHAR_FORMAT, LEAVE_ROOM) + json.toString();
    }

    public static int leaveRoomResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getInt("RequestID");
    }

    public static String createRoomRequest(String roomName, String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("RoomName", roomName);
        json.put("Username", username);

        return String.format(TWO_CHAR_FORMAT, CREATE_ROOM) + json.toString();
    }

    public static int createRoomResponse(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return json.getInt("RoomID");
    }

    public static String postMsgRequest(String username, Message msg) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("Username", username);

        if (msg.roomID == 0) throw new AssertionError();
        json.put("RoomID", msg.roomID);
        json.put("RoomID", msg.text);

        return String.format(TWO_CHAR_FORMAT, MESSAGE) + json.toString();
    }

    public static Message messageRecieved(String rawString) throws JSONException {
        JSONObject json = makeProperJsonObject(rawString);
        return  new Message(json.getJSONObject("Msg"));
    }

    public static String messageReadRequest(String username, int roomID, int msgID) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("Username", username);
        json.put("RoomID", roomID);
        json.put("MsgID", msgID);

        return String.format(TWO_CHAR_FORMAT, MSG_READ) + json.toString();
    }

    public static String signoutRequest(String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("Username", username);

        return String.format(TWO_CHAR_FORMAT, SIGN_OUT) + json.toString();
    }

    public static String deleteUserRequest(String username) throws JSONException {
        JSONObject json = constructJSONRequest();
        json.put("Username", username);

        return String.format(TWO_CHAR_FORMAT, DELETE_USER) + json.toString();
    }



    public static void readErrorResponse(String rawString) {
        try {
            JSONObject json = makeProperJsonObject(rawString);
            String error = json.getString("ErrorCause");
            throw new RuntimeException(error);
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
}
