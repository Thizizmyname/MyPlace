package com.myplace.myplace;

import android.content.SharedPreferences;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.Toast;

import java.util.ArrayList;

import static com.myplace.myplace.LoginActivity.LOGIN_PREFS;

public class MessageActivity extends AppCompatActivity {
    private Toast messageEmptyToast = null;
    RoomDbHelper roomDB = null;

    //TEST FOR INCOMING AND OUTGOING
    int a = 0;
    String name = "Joel";

    //TEST FOR INCOMING AND OUTGOING
    private int mod(int x, int y)
    {
        int result = x % y;
        if (result < 0)
            result += y;
        return result;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_message);
        Room room = getIntent().getParcelableExtra("Room");

        final String roomName = room.getName();

        //noinspection ConstantConditions
        getSupportActionBar().setTitle(roomName);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);

        roomDB = new RoomDbHelper(this);

        ArrayList<Message> messageList = roomDB.getMessages(roomName);
        final MessageAdapter messageAdapter = new MessageAdapter(this, messageList);

        // Finds the listview and specifies the adapter to use
        ListView listMessages = (ListView) findViewById(R.id.listMessages);
        listMessages.setAdapter(messageAdapter);

        final EditText message = (EditText) findViewById(R.id.editMsg);

        final ImageButton btnSend = (ImageButton) findViewById(R.id.btnSendMsg);
        btnSend.setEnabled(false);

        message.addTextChangedListener(onTextChanged);


        btnSend.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {

                // Check if message is empty
                String messageString = message.getText().toString();
                if (messageString.matches("")) {
                    if (messageEmptyToast != null) messageEmptyToast.cancel();
                    messageEmptyToast = Toast.makeText(MessageActivity.this, R.string.message_empty, Toast.LENGTH_SHORT);
                    messageEmptyToast.show();
                    return;
                }

                //TEST FOR INCOMING AND OUTGOING
                if(mod(a, 2) == 0) {
                    String username = "N/A";
                    SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
                    name = loginInfo.getString("username", username);
                } else {
                    name = "Joel";
                }

                Message newMessage = new Message(name, message.getText().toString());
                messageAdapter.add(newMessage);

                roomDB.addMessage(roomName, newMessage);
                MainActivity.roomAdapter.notifyDataSetChanged();

                //TEST FOR INCOMING AND OUTGOING
                ++a; //TEST

                message.setText(null); // Reset input field
            }
        });
    }

    private TextWatcher onTextChanged = new TextWatcher() {
        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {

        }

        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {

        }

        @Override
        public void afterTextChanged(Editable s) {
            ImageButton btnSend = (ImageButton) findViewById(R.id.btnSendMsg);
            if(s.length() > 0) {
                btnSend.setEnabled(true);
            } else {
                btnSend.setEnabled(false);
            }
        }
    };

    @Override
    public void onBackPressed() {
        super.onBackPressed();

        overridePendingTransition(R.anim.push_right_in, R.anim.push_right_out);
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        int id = item.getItemId();

        if (id == android.R.id.home) {
            onBackPressed();  return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
