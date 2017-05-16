package com.myplace.myplace;

import android.app.ProgressDialog;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.ServiceConnection;
import android.os.Bundle;
import android.os.IBinder;
import android.os.StrictMode;
import android.support.v4.content.LocalBroadcastManager;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.util.Patterns;
import android.view.View;
import android.view.WindowManager;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import com.myplace.myplace.services.ConnectionService;
import com.myplace.myplace.services.LoginBroadcastReceiver;

import org.json.JSONException;

import butterknife.ButterKnife;
import butterknife.Bind;

public class SignupActivity extends AppCompatActivity {
    private static final String TAG = "SignupActivity";
    ConnectionService mService;
    boolean mBound = false;
    private String username;
    private volatile Boolean signupAccepted;
    private ProgressDialog progressDialog; // = new ProgressDialog(SignupActivity.this);

    @Bind(R.id.sign_retype) EditText _passRetype;
    @Bind(R.id.sign_username) EditText _userSign;
    @Bind(R.id.sign_password) EditText _passSign;
    @Bind(R.id.sign_btn) Button _btnSign;
    @Bind(R.id.link_login) TextView _linkLogin;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_HIDDEN);
        setContentView(R.layout.activity_signup);
        ButterKnife.bind(this);
        signupAccepted = null;
        progressDialog = new ProgressDialog(SignupActivity.this);

        _btnSign.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                signUp();
            }
        });

        _linkLogin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                onBackPressed();
            }
        });
    }

    public void signUp(){
        Log.d(TAG, "SignUp");

        if(!validate()) {
            onSignUpFailed();
            return;
        }

        _btnSign.setEnabled(false);

        //final ProgressDialog progressDialog = new ProgressDialog(SignupActivity.this);
        progressDialog.setIndeterminate(true);
        progressDialog.setMessage("Creating Account...");
        progressDialog.show();

        username = _userSign.getText().toString();
        final String password = _passSign.getText().toString();

        try {
            mService.sendMessage(JSONParser.signupRequest(username, password));

        } catch (JSONException e) {
            e.printStackTrace();
        }

//        new android.os.Handler().postDelayed(
//                new Runnable() {
//                    public void run() {
//                        try {
//                            mService.sendMessage(JSONParser.signupRequest(username, password));
//
//                        while (signupAccepted == null) {
//                            Thread.sleep(100);
//                        }
//                        if (signupAccepted) {
//                            onSignUpSuccess(username);
//                        }
//                        else {
//                            onSignUpFailed();
//                        }
//                        progressDialog.dismiss();
//                        } catch (InterruptedException | JSONException e) {
//                            e.printStackTrace();
//                        }
//                    }
//                }, 3000);
    }

    public void onSignUpSuccess(String username) {
        _btnSign.setEnabled(true);
        Intent result = new Intent();
        result.putExtra("username", username);
        setResult(RESULT_OK, result);
        finish();
    }

    public void onSignUpFailed() {
        Toast.makeText(getBaseContext(), "Login Failed",  Toast.LENGTH_LONG).show();

        _btnSign.setEnabled(true);
    }

    public boolean validate() {
        Boolean valid = true;

        String username = _userSign.getText().toString();
        String password = _passSign.getText().toString();
        String reType = _passRetype.getText().toString();

        if(!(reType.equals(password))){
            _passRetype.setError(getResources().getString(R.string.error_incorrect_password));
            valid = false;
        }
        else {
            _passRetype.setError(null);
        }

        if(username.isEmpty() || username.length() <= 3) {
            _userSign.setError(getResources().getString(R.string.error_incorrect_username));
            valid = false;
        }
        else {
            _userSign.setError(null);
        }

        if(password.isEmpty() || password.length() <= 5){
            _passSign.setError(getResources().getString(R.string.error_incorrect_password));
            valid = false;
        }
        else {
            _passSign.setError(null);
        }

        return valid;
    }

    @Override
    public void onBackPressed() {
        super.onBackPressed();
        overridePendingTransition(R.anim.push_left_in, R.anim.push_left_out);
    }


    @Override
    protected void onStart() {
        super.onStart();
        // Bind to LocalService
        Log.d("SignupActivity", "Activity onStart!");
        Intent intent = new Intent(this, ConnectionService.class);
        bindService(intent, mTConnection, Context.BIND_AUTO_CREATE);
    }

    @Override
    protected void onStop() {
        super.onStop();
        // Unbind from the service
        if (mBound) {
            unbindService(mTConnection);
            mBound = false;
        }
    }
    @Override
    protected void onResume() {
        super.onResume();
        // Register to receive messages.
        // We are registering an observer (mMessageReceiver) to receive Intents
        // with actions named "custom-event-name".
        LocalBroadcastManager.getInstance(this).registerReceiver(loginBroadcastReceiver,
                new IntentFilter(ConnectionService.BROADCAST_TAG));
    }

    @Override
    protected void onPause() {
        super.onPause();
        LocalBroadcastManager.getInstance(this).unregisterReceiver(loginBroadcastReceiver);
    }

    private LoginBroadcastReceiver loginBroadcastReceiver = new LoginBroadcastReceiver() {
        @Override
        public void handleBooleanResponse(boolean serverResponse) {
            Log.e("Signup Activity", "Response Received: " + serverResponse);
            //signupAccepted = serverResponse;
            progressDialog.dismiss();
            if (serverResponse) {
                onSignUpSuccess(username);
            } else {
                onSignUpFailed();
            }
        }
    };

    /** Defines callbacks for service binding, passed to bindService() */
    private ServiceConnection mTConnection = new ServiceConnection() {

        @Override
        public void onServiceConnected(ComponentName className,
                                       IBinder service) {
            // We've bound to LocalService, cast the IBinder and get LocalService instance
            ConnectionService.ConnectionBinder binder = (ConnectionService.ConnectionBinder) service;
            mService = binder.getService();
            mBound = true;
            //mService.setUpConnection();
        }

        @Override
        public void onServiceDisconnected(ComponentName arg0) {
            mBound = false;
        }
    };
}
