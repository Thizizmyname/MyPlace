package com.myplace.myplace;

import android.app.ProgressDialog;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;

import android.content.Intent;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import butterknife.ButterKnife;
import butterknife.Bind;

public class LoginActivity extends AppCompatActivity {
    private static final String TAG = "LoginActivity";
    private static final int REQUEST_SIGNUP = 0;
    public static final String LOGIN_PREFS = "login_prefs";

    @Bind(R.id.input_username) EditText _username;
    @Bind(R.id.input_password) EditText _password;
    @Bind(R.id.btn_login) Button _login;
    @Bind(R.id.btn_bypass) Button _bypass;
    @Bind(R.id.link_signup) TextView _signup;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);
        ButterKnife.bind(this);
        getSupportActionBar().hide();

        SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
        boolean loggedIn = loginInfo.getBoolean("loggedIn", false);
        _username.setText(String.valueOf(loggedIn));
        if(loggedIn == true){
            Intent startMain = new Intent(getApplicationContext(), MainActivity.class);
            startActivity(startMain);
            finish();
        }

        _bypass.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View v) {
                onLoginSuccess("");
            }
        });

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

    public void login(){
        Log.d(TAG, "login");

        if(!validate()) {
            onLoginFailed();
            return;
        }

        _login.setEnabled(false);

        final ProgressDialog progressDialog = new ProgressDialog(LoginActivity.this);
        progressDialog.setIndeterminate(true);
        progressDialog.setMessage("Authenticating");
        progressDialog.show();

        final String username = _username.getText().toString();
        String password = _password.getText().toString();

        new android.os.Handler().postDelayed(
                new Runnable() {
                    public void run() {
                        onLoginSuccess(username);
                        progressDialog.dismiss();
                    }
                }, 3000);
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        if (requestCode == REQUEST_SIGNUP) {
            if (resultCode == RESULT_OK) {
                String username = data.getStringExtra("username");
                onLoginSuccess(username);
            }
        }
    }

    @Override
    public void onBackPressed() {
        moveTaskToBack(true);
    }

    public void onLoginSuccess(String username) {
        _login.setEnabled(true);
        SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
        SharedPreferences.Editor loginEdit = loginInfo.edit();
        loginEdit.putBoolean("loggedIn", true);
        loginEdit.putString("username", username);
        loginEdit.commit();
        Intent startMain = new Intent(getApplicationContext(), MainActivity.class);
        startActivity(startMain);
        finish();
    }

    public void onLoginFailed() {
        Toast.makeText(getBaseContext(), "Login failed", Toast.LENGTH_LONG).show();

        _login.setEnabled(true);
    }

    public boolean validate() {
        boolean valid = true;

        String username = _username.getText().toString();
        String password = _password.getText().toString();

        if (username.isEmpty() || username.length() <= 3) {
            _username.setError("Username must be atleast 4 characteters");
            valid = false;
        } else {
            _username.setError(null);
        }

        if (password.isEmpty() || password.length() <= 5) {
            _password.setError("Password must be atleast 6 characters");
            valid = false;
        } else {
            _password.setError(null);
        }

        return valid;
    }

}