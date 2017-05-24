package com.myplace.myplace;

import android.app.ProgressDialog;
import android.content.ComponentName;
import android.content.Context;
import android.content.IntentFilter;
import android.content.ServiceConnection;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.os.Handler;
import android.os.IBinder;
import android.support.v4.content.LocalBroadcastManager;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;

import android.content.Intent;
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

import static java.lang.Thread.sleep;

public class LoginActivity extends AppCompatActivity {
    ConnectionService mService;
    boolean mBound = false;
    private static final String TAG = "LoginActivity";
    private static final int REQUEST_SIGNUP = 0;
    public static final String LOGIN_PREFS = "login_prefs";
    private ProgressDialog progressDialog;
    private String username;
    private String password;
    private Handler loginHandler = new Handler();

    @Bind(R.id.input_username) EditText _username;
    @Bind(R.id.input_password) EditText _password;
    @Bind(R.id.btn_login) Button _login;
    @Bind(R.id.link_signup) TextView _signup;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_HIDDEN);
        setContentView(R.layout.activity_login);
        ButterKnife.bind(this);
        getSupportActionBar().hide();

        startService(new Intent(this, ConnectionService.class));

        SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
        boolean loggedIn = loginInfo.getBoolean("loggedIn", false);
        if (loggedIn == true) {
            autoLogin(loginInfo);
        }

        _login.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View v) {
                login();
            }
        });
        _signup.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View v) {
                Intent intent = new Intent(getApplicationContext(), SignupActivity.class);
                startActivityForResult(intent, REQUEST_SIGNUP);
                overridePendingTransition(R.anim.push_left_in, R.anim.push_left_out);
            }
        });
    }

    private void autoLogin(final SharedPreferences loginInfo) {

        _login.setEnabled(false);

        progressDialog = new ProgressDialog(LoginActivity.this);
        progressDialog.setIndeterminate(true);
        progressDialog.setMessage("Authenticating");
        progressDialog.show();
        LocalBroadcastManager.getInstance(getApplicationContext()).registerReceiver(loginBroadcastReceiver,
                new IntentFilter(ConnectionService.BROADCAST_TAG));

        Thread thread = new Thread(new Runnable() {
            @Override
            public void run() {
                while(!mBound){
                    try {
                        sleep(100);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }

                username = loginInfo.getString("username", "No_username_found");
                password = loginInfo.getString("password", "No_password_found");


                loginHandler.postDelayed(failLoginAfterDelay, 10000);

                try {
                    mService.sendMessage(JSONParser.signinRequest(username, password));

                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        });
        thread.start();
    }

    public void login() {
        Log.d(TAG, "login");

        if (!validate()) {
            onLoginFailed();
            return;
        }

        _login.setEnabled(false);

        progressDialog = new ProgressDialog(LoginActivity.this);
        progressDialog.setIndeterminate(true);
        progressDialog.setMessage("Authenticating");
        progressDialog.show();

        username = _username.getText().toString();
        password = _password.getText().toString();

        LocalBroadcastManager.getInstance(this).registerReceiver(loginBroadcastReceiver,
                new IntentFilter(ConnectionService.BROADCAST_TAG));

        loginHandler.postDelayed(failLoginAfterDelay, 10000);

        try {
            mService.sendMessage(JSONParser.signinRequest(username, password));

        } catch (JSONException e) {
            e.printStackTrace();
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        if (requestCode == REQUEST_SIGNUP) {
            if (resultCode == RESULT_OK) {
                username = data.getStringExtra("username");
                onLoginSuccess();
            }
        }
    }

    @Override
    public void onBackPressed() {
        moveTaskToBack(true);
    }

    public void onLoginSuccess() {
        Log.e("onLoginSuccess", "Starting mainActivity");
        _login.setEnabled(true);
        SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
        SharedPreferences.Editor loginEdit = loginInfo.edit();
        loginEdit.putBoolean("loggedIn", true);
        loginEdit.putString("username", username);
        loginEdit.putString("password", password);
        loginEdit.commit();
        Intent startMain = new Intent(getApplicationContext(), MainActivity.class);
        startActivity(startMain);
        finish();
    }

    public void onLoginFailed() {
        Toast.makeText(getBaseContext(), "Login failed", Toast.LENGTH_LONG).show();

        _login.setEnabled(true);
    }

    public void loginConnectionFailed() {
        Toast.makeText(getBaseContext(), "Connection Failed", Toast.LENGTH_LONG).show();

        _login.setEnabled(true);
    }

    public boolean validate() {
        boolean valid = true;

        String valiUsername = _username.getText().toString();
        String valiPassword = _password.getText().toString();

        if (valiUsername.isEmpty() || valiUsername.length() <= 3) {
            _username.setError(getResources().getString(R.string.error_incorrect_username));
            valid = false;
        } else {
            _username.setError(null);
        }

        if (valiPassword.isEmpty() || valiPassword.length() <= 5) {
            _password.setError(getResources().getString(R.string.error_incorrect_password));
            valid = false;
        } else {
            _password.setError(null);
        }
        return valid;
    }

    @Override
    protected void onStart() {
        super.onStart();
        // Bind to LocalService
        Log.d("Main_Activity", "I'm in onStart!");
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

    private Runnable failLoginAfterDelay = new Runnable() {
        @Override
        public void run() {
            progressDialog.dismiss();
            LocalBroadcastManager.getInstance(getApplicationContext()).unregisterReceiver(loginBroadcastReceiver);
            loginConnectionFailed();
        }
    };

    private LoginBroadcastReceiver loginBroadcastReceiver = new LoginBroadcastReceiver() {
        public void handleBooleanResponse(boolean serverResponse) {
            Log.d("Login Activity", "Response Received: " + serverResponse);
            progressDialog.dismiss();
            loginHandler.removeCallbacks(failLoginAfterDelay);
            LocalBroadcastManager.getInstance(getApplicationContext()).unregisterReceiver(loginBroadcastReceiver);
            if (serverResponse) {
                onLoginSuccess();
            } else {
                onLoginFailed();
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
        }

        @Override
        public void onServiceDisconnected(ComponentName arg0) {
            mBound = false;
        }
    };

}