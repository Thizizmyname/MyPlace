package com.myplace.myplace;

import android.content.ContentValues;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.Toast;

import java.util.ArrayList;

import static com.myplace.myplace.MainActivity.roomDB;

public class MessageActivity extends AppCompatActivity {
    private Toast messageEmptyToast = null;

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

    public ArrayList<Message> getMessages(String roomName){
        String query = "SELECT * FROM "+roomName;
        Cursor c = roomDB.rawQuery(query,null);
        ArrayList<Message> list=new ArrayList<>();
        c.moveToFirst();

        while(c.moveToNext()){
            String name = c.getString(c.getColumnIndex("name"));

            String message = c.getString(c.getColumnIndex("message"));

            String date = c.getString(c.getColumnIndex("date"));

            Message newMessage = new Message(name, message, date);
            list.add(newMessage);
        }

        c.close();
        return list;
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

        // Creates an array containing the messages and an adapter
        // final ArrayList<Message> messageArray = new ArrayList<>();

        ArrayList<Message> messageList = getMessages(roomName);
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

                ContentValues insertValues = new ContentValues();
                insertValues.put("name", newMessage.getName());
                insertValues.put("message", newMessage.text);
                insertValues.put("date", newMessage.date);
                roomDB.insert(roomName, null, insertValues);

                //roomDB.execSQL("INSERT INTO "+roomName+" VALUES('" + newMessage.getName() + "','" + newMessage.text + "','" + newMessage.date + "');");

                //room.messageList.add(newMessage);

                //TEST FOR INCOMING AND OUTGOING
                ++a; //TEST

                message.setText(null); // Reset input field
            }
        });
    }
}
