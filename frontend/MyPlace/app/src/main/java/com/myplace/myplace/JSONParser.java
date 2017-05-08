package com.myplace.myplace;


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

    public static String getRoomrequest(String userName) throws JSONException {
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
            int _id = r.getInt("RoomID");
            String _name = r.getString("RoomName");
            int _latest = r.getInt("LatestReadMsgID");

            JSONObject ms = r.getJSONObject("LatestMsg");
            //Message msg = new Message();


        }
        return null;
    }


}
