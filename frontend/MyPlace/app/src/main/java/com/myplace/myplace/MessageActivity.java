package com.myplace.myplace;

import android.content.BroadcastReceiver;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.ServiceConnection;
import android.content.SharedPreferences;
import android.os.IBinder;
import android.support.v4.content.LocalBroadcastManager;
import android.support.v4.widget.SwipeRefreshLayout;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.Toast;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;
import com.myplace.myplace.services.ConnectionService;
import com.myplace.myplace.services.MainBroadcastReceiver;

import org.json.JSONException;

import java.util.ArrayList;

import static com.myplace.myplace.LoginActivity.LOGIN_PREFS;

public class MessageActivity extends AppCompatActivity {

    private static final String EMPTY_STRING = "";

    MessageAdapter messageAdapter;
    private Toast messageEmptyToast = null;
    RoomDbHelper roomDB = null;

    ConnectionService mService;
    boolean mBound = false;
    private int roomID;
    private SwipeRefreshLayout swipeContainer;


    // Our handler for received Intents. This will be called whenever an Intent
    // with an action named "custom-event-name" is broadcasted.
    private MainBroadcastReceiver mMessageReceiver = new MainBroadcastReceiver() {
        @Override
        public void handleCreatedRoomInActivity(Room room) {

        }

        @Override
        public void handleNewMessageInActivity(Message msg) {
            if (roomID == msg.getRoomID()) {
                messageAdapter.add(msg);
            }
        }

        @Override
        public void handleJoinedRoomInActivity(RoomInfo roominfo) {

        }

        @Override
        public void handleOlderMessagesInActivity(ArrayList<Message> messages) {
            messageAdapter.updateData(roomDB.getMessages(roomID));
            swipeContainer.setRefreshing(false);
        }
    };

    @Override
    protected void onStart() {
        super.onStart();
        // Bind to LocalService
        Log.d("MessageActivity", "I'm in onStart!");
        Intent intent = new Intent(this, ConnectionService.class);
        bindService(intent, mTConnection, Context.BIND_AUTO_CREATE);
    }

    @Override
    protected void onStop() {
        super.onStop();
        // Unbind from the service
        if (mBound) {
            Log.d("MessageActivity", "Stopping event");
            unbindService(mTConnection);
            mBound = false;
        }
    }

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

    @Override
    protected void onResume() {
        super.onResume();
        // Register to receive messages.
        // We are registering an observer (mMessageReceiver) to receive Intents
        // with actions named "custom-event-name".
        LocalBroadcastManager.getInstance(this).registerReceiver(mMessageReceiver,
                new IntentFilter(ConnectionService.BROADCAST_TAG));
    }

    @Override
    protected void onPause() {
        super.onPause();
        LocalBroadcastManager.getInstance(this).unregisterReceiver(mMessageReceiver);
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_message);

        final String roomName = getIntent().getExtras().getString(MainActivity.ROOM_NAME);
        roomID = getIntent().getExtras().getInt("roomID");

        //noinspection ConstantConditions
        getSupportActionBar().setTitle(roomName);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);

        roomDB = new RoomDbHelper(this);

        // Initialize swipecontainer
        swipeContainer = (SwipeRefreshLayout) findViewById(R.id.swipeContainer);
        swipeContainer.setOnRefreshListener(new SwipeRefreshLayout.OnRefreshListener() {
            @Override
            public void onRefresh() {
                try {
                    if (!messageAdapter.isEmpty()) {
                        Message oldestMessage = messageAdapter.getItem(0);
                        mService.sendMessage(JSONParser.getOlderMsgsRequest(roomID, oldestMessage.getId()));
                    }
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        });

        ArrayList<Message> messageList = roomDB.getMessages(roomID);
        messageAdapter = new MessageAdapter(this, messageList);

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
                if (messageString.matches(EMPTY_STRING)) {
                    if (messageEmptyToast != null) messageEmptyToast.cancel();
                    messageEmptyToast = Toast.makeText(MessageActivity.this, R.string.message_empty, Toast.LENGTH_SHORT);
                    messageEmptyToast.show();
                    return;
                }

                SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
                final String username = loginInfo.getString("username", MainActivity.NO_USERNAME_FOUND);

                long timestamp = System.currentTimeMillis();

                Message newMessage = new Message(roomID, username, message.getText().toString(), timestamp);

                try {
                    mService.sendMessage(JSONParser.postMsgRequest(newMessage));

                } catch (JSONException e) {
                    e.printStackTrace();
                }
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
