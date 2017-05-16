package com.myplace.myplace.services;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;

import com.myplace.myplace.JSONParser;
import com.myplace.myplace.RoomDbHelper;
import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.RequestTypes;

import org.json.JSONException;

/**
 * Created by alexis on 2017-05-15.
 */

public abstract class MainBroadcastReceiver extends BroadcastReceiver {

    Context context;
    RoomDbHelper roomDB;

    public MainBroadcastReceiver(Context context) {
        this.context = context;
        this.roomDB = new RoomDbHelper(context);
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        // Get extra data included in the Intent
        String serverMessage = intent.getStringExtra(ConnectionService.REPLY_PACKAGE);

        int i = JSONParser.determineJSONType(serverMessage);

        try {

            switch (i) {
                case RequestTypes.SIGN_UP:
                case RequestTypes.SIGN_IN:
                    Log.e("MainBroadcastReceiver", "Signup/Signin responses in Main/MessageActivity");
                    break;
                case RequestTypes.GET_ROOMS:

                    break;
                case RequestTypes.GET_USERS:

                    break;
                case RequestTypes.GET_OLDER:

                    break;
                case RequestTypes.GET_NEWER:

                    break;
                case RequestTypes.JOIN_ROOM:

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

    private void newMessageReceived(final Message message) {
        final int roomID = message.getRoomID();

        roomDB.addMessage(Integer.toString(roomID), message);
        handleNewMessageInActivity(message);

    }


    public abstract void handleNewMessageInActivity(final Message msg);


}
