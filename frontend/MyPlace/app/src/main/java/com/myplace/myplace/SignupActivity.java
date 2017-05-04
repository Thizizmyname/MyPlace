package com.myplace.myplace;

import android.app.ProgressDialog;
import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.util.Patterns;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import butterknife.ButterKnife;
import butterknife.Bind;

public class SignupActivity extends AppCompatActivity {
    private static final String TAG = "SignupActivity";

    @Bind(R.id.sign_email) EditText _emailSign;
    @Bind(R.id.sign_username) EditText _userSign;
    @Bind(R.id.sign_password) EditText _passSign;
    @Bind(R.id.sign_btn) Button _btnSign;
    @Bind(R.id.link_login) TextView _linkLogin;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_signup);
        ButterKnife.bind(this);

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

        final ProgressDialog progressDialog = new ProgressDialog(SignupActivity.this);
        progressDialog.setIndeterminate(true);
        progressDialog.setMessage("Creating Account...");
        progressDialog.show();

        final String username = _userSign.getText().toString();
        String password = _passSign.getText().toString();

        new android.os.Handler().postDelayed(
                new Runnable() {
                    public void run() {
                        onSignUpSuccess(username);
                        progressDialog.dismiss();
                    }
                }, 3000);
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


        if(username.isEmpty() || username.length() <= 3) {
            _userSign.setText("Username must be longer than 3 characters");
            valid = false;
        }
        else {
            _userSign.setError(null);
        }

        if(password.isEmpty() || password.length() <= 5){
            _passSign.setText("Password must be atleast 6 characters");
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
}
