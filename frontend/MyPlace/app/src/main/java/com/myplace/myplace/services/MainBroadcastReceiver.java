package com.myplace.myplace.services;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;

import com.myplace.myplace.JSONParser;
import com.myplace.myplace.RoomDbHelper;
import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.RequestTypes;
import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;

import org.json.JSONException;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-05-15.
 */

public abstract class MainBroadcastReceiver extends BroadcastReceiver {


    RoomDbHelper roomDB = null;


    @Override
    public void onReceive(Context context, Intent intent) {

        if (roomDB == null) {
            roomDB = new RoomDbHelper(context);
        }
        // Get extra data included in the Intent
        String serverMessage = intent.getStringExtra(ConnectionService.REPLY_PACKAGE);

        Log.d("MainBroadcastReceiver", "Received: " + serverMessage);
        int i = JSONParser.determineJSONType(serverMessage);
        Log.d("MainBroadcastReceiver", "JSONType: " + i);

        try {

            switch (i) {
                case RequestTypes.SIGN_UP:
                case RequestTypes.SIGN_IN:
                    Log.e("MainBroadcastReceiver", "Signup/Signin responses in Main/MessageActivity");
                    break;
                case RequestTypes.GET_ROOMS:
                    final ArrayList<RoomInfo> roomResponse = JSONParser.getRoomResponse(serverMessage);
                    handleRoomList(roomResponse);
                    break;
                case RequestTypes.GET_USERS:

                    break;
                case RequestTypes.GET_OLDER:
                    ArrayList<Message> messages = JSONParser.getOlderMsgsResponse(serverMessage);
                    handleOlderMessages(messages);
                    break;
                case RequestTypes.GET_NEWER:

                    break;
                case RequestTypes.CREATE_ROOM:
                    Room room = JSONParser.createRoomResponse(serverMessage);
                    handleCreatedRoom(room);
                    break;
                case RequestTypes.JOIN_ROOM:
                    RoomInfo joinroom = JSONParser.joinRoomResponse(serverMessage);
                    handleJoinedRoom(joinroom);
                    break;
                case RequestTypes.LEAVE_ROOM:

                    break;
                case RequestTypes.MESSAGE:
                    final Message message = JSONParser.messageRecieved(serverMessage);
                    newMessageReceived(message);
                    break;
                case RequestTypes.MSG_READ:

                    break;
                case RequestTypes.SIGN_OUT:

                    break;
                case RequestTypes.DELETE_USER:

                    break;
                case RequestTypes.ERROR_TYPE:
                    JSONParser.throwErrorResponse(serverMessage);
                    break;
                default:
                    // Undefined server response
                    throw new RuntimeException(serverMessage);

            }
        } catch (JSONException e) {
            e.printStackTrace();
        }
    }

    protected void handleRoomList(ArrayList<RoomInfo> roomResponse) {
        throw new RuntimeException("No implementation");
    }


    private void newMessageReceived(final Message message) {

        roomDB.addMessage(message.getRoomID(), message);
        handleNewMessageInActivity(message);

    }

    private void handleCreatedRoom(final Room room) {
        roomDB.createRoomTable(room);
        handleCreatedRoomInActivity(room);
    }

    private void handleJoinedRoom(final RoomInfo room) {
        roomDB.createRoomTable(room.getRoom());
        handleJoinedRoomInActivity(room);
    }

    private void handleOlderMessages(final ArrayList<Message> messages) {
        for (Message message : messages) {
            roomDB.addMessage(message.getRoomID(), message);
        }
        handleOlderMessagesInActivity(messages);

    }

    public void handleJoinedRoomInActivity(final RoomInfo roominfo) {}

    public void handleOlderMessagesInActivity(final ArrayList<Message> messages){}

    public void handleCreatedRoomInActivity(final Room room){}

    public abstract void handleNewMessageInActivity(final Message msg);


}
