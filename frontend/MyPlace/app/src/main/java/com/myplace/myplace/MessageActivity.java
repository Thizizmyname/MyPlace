package com.myplace.myplace;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.TextView;

import java.util.ArrayList;
import java.util.List;

public class MessageActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_message);

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
                final String name = "Anders";

                Message newMessage = new Message(name, message.getText().toString());
                messageAdapter.add(newMessage);

                message.setText(null); // Reset input field
            }
        });
    }
}
