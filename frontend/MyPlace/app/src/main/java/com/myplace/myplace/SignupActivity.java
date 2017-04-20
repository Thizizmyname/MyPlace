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
                Intent intent = new Intent(getApplicationContext(), LoginActivity.class);
                startActivity(intent);
                finish();
                overridePendingTransition(R.anim.push_left_in, R.anim.push_left_out);
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

        String user = _userSign.getText().toString();
        String email = _emailSign.getText().toString();
        String password = _passSign.getText().toString();

        new android.os.Handler().postDelayed(
                new Runnable() {
                    public void run() {
                        onSignUpSucces();
                        progressDialog.dismiss();
                    }
                }, 3000);
    }

    public void onSignUpSucces() {
        _btnSign.setEnabled(true);
        setResult(RESULT_OK, null);
        finish();
    }

    public void onSignUpFailed() {
        Toast.makeText(getBaseContext(), "Login Failed",  Toast.LENGTH_LONG).show();

        _btnSign.setEnabled(true);
    }

    public boolean validate() {
        Boolean valid = true;

        String email = _emailSign.getText().toString();
        String username = _userSign.getText().toString();
        String password = _passSign.getText().toString();

        if (email.isEmpty() || !Patterns.EMAIL_ADDRESS.matcher(email).matches()) {
            _emailSign.setText("Enter a vaild Email");
            valid = false;
        }
        else{
            _emailSign.setError(null);
        }

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
}
