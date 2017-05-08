package com.myplace.myplace;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.Toast;

import com.myplace.myplace.models.Message;

import java.util.ArrayList;

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

        final String roomName = getIntent().getExtras().getString("RoomName");

        //noinspection ConstantConditions
        getSupportActionBar().setTitle(roomName);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);

        roomDB = new RoomDbHelper(this);

        ArrayList<Message> messageList = roomDB.getMessages(roomName);
        final MessageAdapter messageAdapter = new MessageAdapter(this, messageList);

        // Finds the listview and specifies the adapter to use
        ListView listMessages = (ListView) findViewById(R.id.listMessages);
        listMessages.setAdapter(messageAdapter);

        final ImageButton btnSend = (ImageButton) findViewById(R.id.btnSendMsg);

        btnSend.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                final EditText message = (EditText) findViewById(R.id.editMsg);
                String messageString = message.getText().toString();

                // Check if message is empty
                if (messageString.matches("")) {
                    if (messageEmptyToast != null) messageEmptyToast.cancel();
                    messageEmptyToast = Toast.makeText(MessageActivity.this, R.string.message_empty, Toast.LENGTH_SHORT);
                    messageEmptyToast.show();
                    return;
                }

                //TEST FOR INCOMING AND OUTGOING
                if(mod(a, 2) == 0) {
                    name = "Anders";
                } else {
                    name = "Joel";
                }

                Message newMessage = new Message(name, message.getText().toString());
                messageAdapter.add(newMessage);

                roomDB.addMessage(roomName, newMessage);
                MainActivity.adapter.notifyDataSetChanged();

                //TEST FOR INCOMING AND OUTGOING
                ++a; //TEST

                message.setText(null); // Reset input field
            }
        });
    }
}
