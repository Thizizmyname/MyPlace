package com.myplace.myplace.services;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;

import com.myplace.myplace.JSONParser;
import com.myplace.myplace.models.RequestTypes;

import org.json.JSONException;

/**
 * Created by alexis on 2017-05-16.
 */

public abstract class LoginBroadcastReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        // Get extra data included in the Intent
        String serverMessage = intent.getStringExtra(ConnectionService.REPLY_PACKAGE);

        int i = JSONParser.determineJSONType(serverMessage);

        try {

            switch (i) {
                case RequestTypes.SIGN_UP:
                    final Boolean signupResponse = JSONParser.signupResponse(serverMessage);
                    handleBooleanResponse(signupResponse);
                    break;
                case RequestTypes.SIGN_IN:
                    final Boolean signinResponse = JSONParser.signinResponse(serverMessage);
                    handleBooleanResponse(signinResponse);
                    break;
                case RequestTypes.ERROR_TYPE:
                    JSONParser.throwErrorResponse(serverMessage);
                    break;
                default:
                    // Don't handle other responses response
                    break;

            }
        } catch (JSONException e) {
            e.printStackTrace();
        }
    }

    public abstract void handleBooleanResponse(boolean serverResponse);

}
