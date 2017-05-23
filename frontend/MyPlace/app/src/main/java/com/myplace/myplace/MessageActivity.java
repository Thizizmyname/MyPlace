package com.myplace.myplace;

import android.content.BroadcastReceiver;
import android.content.ComponentName;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.ServiceConnection;
import android.content.SharedPreferences;
import android.os.IBinder;
import android.support.v4.content.LocalBroadcastManager;
import android.support.v4.widget.SwipeRefreshLayout;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;
import com.myplace.myplace.services.ConnectionService;
import com.myplace.myplace.services.MainBroadcastReceiver;

import org.json.JSONException;
import org.w3c.dom.Text;

import java.util.ArrayList;

import static com.myplace.myplace.LoginActivity.LOGIN_PREFS;
import static java.lang.Thread.sleep;

public class MessageActivity extends AppCompatActivity {
    private String TAG = "MessageActivity";
    final Context context = this;
    private static final String EMPTY_STRING = "";
    private static final int FIRST_MSGID_IN_CONVERSATION = 0;

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
        public void handleNewMessageInActivity(Message msg) {
            if (roomID == msg.getRoomID()) {
                messageAdapter.add(msg);
            }
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
        Log.d(TAG, "I'm in onStart!");
        Intent intent = new Intent(this, ConnectionService.class);
        bindService(intent, mTConnection, Context.BIND_AUTO_CREATE);
    }

    @Override
    protected void onStop() {
        super.onStop();
        // Unbind from the service
        if (mBound) {
            Log.d(TAG, "Stopping event");
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
        getOlderIfNeeded(messageList);
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
                onSendButtonClick(message);
            }
        });
    }

    private void onSendButtonClick(EditText message) {
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

    public void onShowRoomId() {
        LayoutInflater inflater = LayoutInflater.from(context);
        View dialogView = inflater.inflate(R.layout.dialog_show_roomid, null);
        AlertDialog.Builder builder = new AlertDialog.Builder(context);

        builder.setView(dialogView);
        TextView roomIdTextView = (TextView) dialogView.findViewById(R.id.room_id);
        roomIdTextView.setText(Integer.toString(roomID));

        builder.setPositiveButton("Ok", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
            // Cancel
            }
        });
        builder.create();
        builder.show();
    }

    private void getOlderIfNeeded(final ArrayList<Message> msgList) {

        int listLength = msgList.size();
        if (listLength == 1 && msgList.get(0).getRoomID() != FIRST_MSGID_IN_CONVERSATION) {

            Thread thread = new Thread(new Runnable() {
                @Override
                public void run() {
                    while (!mBound) {
                        Log.d("MainActivity", "Waiting for mBound");
                        try {
                            sleep(500);
                        } catch (InterruptedException e) {
                            e.printStackTrace();
                        }
                    }
                    try {
                        mService.sendMessage(JSONParser.getOlderMsgsRequest(roomID, msgList.get(0).getId()));
                    } catch (JSONException e) {
                        Log.d("MainActivity", "Get room request error");
                        e.printStackTrace();
                    }
                }
            });

            thread.start();
        }
    }

    @Override
    public void onBackPressed() {
        super.onBackPressed();

        overridePendingTransition(R.anim.push_right_in, R.anim.push_right_out);
    }

    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.message_menu, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        int id = item.getItemId();

        if (id == android.R.id.home) {
            onBackPressed();  return true;
        } else if (id == R.id.show_roomid) {
            onShowRoomId();
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
