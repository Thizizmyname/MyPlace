package com.myplace.myplace;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;

import java.util.ArrayList;

public class MessageActivity extends AppCompatActivity {


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
        //noinspection ConstantConditions
        getSupportActionBar().setTitle(room.getName());
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);


        // Creates an array containing the messages and an adapter
        final ArrayList<Message> messageArray = new ArrayList<>();
        final MessageAdapter messageAdapter = new MessageAdapter(this, messageArray);

        // Finds the listview and specifies the adapter to use
        ListView listMessages = (ListView) findViewById(R.id.listMessages);
        listMessages.setAdapter(messageAdapter);

        final ImageButton btnSend = (ImageButton) findViewById(R.id.btnSendMsg);

        btnSend.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                final EditText message = (EditText) findViewById(R.id.editMsg);

                //TEST FOR INCOMING AND OUTGOING
                if(mod(a, 2) == 0) {
                    name = "Anders";
                } else {
                    name = "Joel";
                }

                Message newMessage = new Message(name, message.getText().toString());
                messageAdapter.add(newMessage);
                //room.messageList.add(newMessage);

                //TEST FOR INCOMING AND OUTGOING
                ++a; //TEST

                message.setText(null); // Reset input field
            }
        });
    }
}
