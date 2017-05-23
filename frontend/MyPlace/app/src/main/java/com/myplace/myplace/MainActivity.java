package com.myplace.myplace;

import android.content.BroadcastReceiver;
import android.content.ComponentName;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.ServiceConnection;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.os.IBinder;
import android.support.v4.content.LocalBroadcastManager;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.text.InputType;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.WindowManager;
import android.widget.AdapterView;
import android.widget.EditText;
import android.widget.ListView;

import com.myplace.myplace.models.Message;
import com.myplace.myplace.models.Room;
import java.util.ArrayList;

import com.getbase.floatingactionbutton.FloatingActionButton;
import com.getbase.floatingactionbutton.FloatingActionsMenu;
import com.myplace.myplace.models.RoomInfo;
import com.myplace.myplace.services.ConnectionService;
import com.myplace.myplace.services.MainBroadcastReceiver;


import org.json.JSONException;

import static com.myplace.myplace.LoginActivity.LOGIN_PREFS;
import static com.myplace.myplace.R.id.action_create;
import static com.myplace.myplace.R.id.action_join;
import static java.lang.Thread.sleep;

public class MainActivity extends AppCompatActivity {

    protected static final String ROOM_NAME = "RoomName";
    protected static final String NO_USERNAME_FOUND = "N/A";

    final Context context = this;
    ConnectionService mService;
    boolean mBound = false;

    private static String username;
    FloatingActionsMenu actionMenu;
    ListView listView;
    ArrayList<RoomInfo> roomList = null;
    public static RoomAdapter roomAdapter = null;

    //Defines the database
    public RoomDbHelper roomDB = null;

    // Our handler for received Intents. This will be called whenever an Intent
    // with an action named "custom-event-name" is broadcasted.
    private MainBroadcastReceiver mMessageReceiver = new MainBroadcastReceiver() {
        @Override
        public void handleRoomListInActivity(ArrayList<RoomInfo> roomlist) {
            ArrayList<RoomInfo> updatedRoomList = roomDB.getRoomList();
            roomAdapter.updateData(updatedRoomList);
            roomAdapter.notifyDataSetChanged();
        }

        @Override
        public void handleJoinedRoomInActivity(RoomInfo roominfo) {
            roomAdapter.add(roominfo);
            roomAdapter.notifyDataSetChanged();
        }

        @Override
        public void handleCreatedRoomInActivity(RoomInfo roominfo) {
            roomAdapter.add(roominfo);
            roomAdapter.notifyDataSetChanged();
        }

        @Override
        public void handleNewMessageInActivity(Message msg) {
            ArrayList<RoomInfo> updatedRoomList = roomDB.getRoomList();
            roomAdapter.updateData(updatedRoomList);
            roomAdapter.notifyDataSetChanged();
        }

        @Override
        public void handleLeaveRoomInActivity(){

        }

        @Override
        public void handleLogoutInActivity() {
            roomDB.dropAllTables();
            logout();

            Intent startLogin = new Intent(getApplicationContext(), LoginActivity.class);
            startActivity(startLogin);
            finish();
        }

        @Override
        public void handleMessageReadInActivity() {
            MainActivity.roomAdapter.notifyDataSetChanged();
        }
    };



    @Override
    protected void onStart() {
        super.onStart();
        // Bind to LocalService
        Log.d("Main_Activity", "I'm in onStart!");
        Intent intent = new Intent(this, ConnectionService.class);
        bindService(intent, mTConnection, Context.BIND_AUTO_CREATE);

        Thread thread = new Thread(new Runnable() {
            @Override
            public void run() {
                while(!mBound) {
                    Log.d("MainActivity", "Waiting for mBound");
                    try {
                        sleep(1000);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
                try {
                    mService.sendMessage(JSONParser.getRoomRequest(username));
                } catch (JSONException e) {
                    Log.d("MainActivity", "Get room request error");
                    e.printStackTrace();
                }
            }
        });

        thread.start();
    }

    @Override
    protected void onResume() {
        super.onResume();
        // Register to receive messages.
        // We are registering an observer (mMessageReceiver) to receive Intents
        // with actions named "custom-event-name".
        LocalBroadcastManager.getInstance(this).registerReceiver(mMessageReceiver,
                new IntentFilter(ConnectionService.BROADCAST_TAG));

        ArrayList<RoomInfo> updatedRoomList = roomDB.getRoomList();
        roomAdapter.updateData(updatedRoomList);
        roomAdapter.notifyDataSetChanged();
    }

    @Override
    protected void onPause() {
        super.onPause();
        LocalBroadcastManager.getInstance(this).unregisterReceiver(mMessageReceiver);
    }

    @Override
    protected void onStop() {
        super.onStop();
        // Unbind from the service
        if (mBound) {
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
            //mService.setUpConnection();
        }

        @Override
        public void onServiceDisconnected(ComponentName arg0) {
            mBound = false;
        }
    };

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        Toolbar toolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);
        //noinspection ConstantConditions
        getSupportActionBar().setTitle(getResources().getString(R.string.app_name));

        //Open database
        roomDB = new RoomDbHelper(this);
        roomList = roomDB.getRoomList();
        listView = (ListView) findViewById(R.id.roomList);

        roomAdapter = new RoomAdapter(this, roomList);
        listView.setAdapter(roomAdapter);

        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                onThreadClick(position);
            }
        });

        listView.setOnItemLongClickListener(new AdapterView.OnItemLongClickListener() {
            @Override
            public boolean onItemLongClick(AdapterView<?> parent, View view, final int position, long id) {
                onThreadLongClick(position);
                return true;
            }
        });


        SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
        username = loginInfo.getString("username", NO_USERNAME_FOUND);

        actionMenu = (FloatingActionsMenu) findViewById(R.id.action_menu);

        //OnClick for createRoom
        final FloatingActionButton actionCreate = (FloatingActionButton) findViewById(action_create);

        actionCreate.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                onAddRoomClick(R.string.create_room);
            }
        });


        //OnClick for joinRoom
        final FloatingActionButton actionJoin = (FloatingActionButton) findViewById(action_join);

        actionJoin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                onAddRoomClick(R.string.join_room);
            }
        });
    }



    public void onThreadClick(int position) {
        Intent intent = new Intent(MainActivity.this, MessageActivity.class);
        intent.putExtra(ROOM_NAME, roomList.get(position).getName());
        intent.putExtra("roomID", roomList.get(position).getRoomID());
        startActivity(intent);
        overridePendingTransition(R.anim.push_left_in, R.anim.push_left_out);
    }

    public void onThreadLongClick(final int position) {
        final int roomID = roomList.get(position).getRoomID();
        final String roomName = roomList.get(position).getName();

        AlertDialog.Builder builder = new AlertDialog.Builder(context);
        builder.setMessage("Do you want to leave "+roomName+"?");
        builder.setPositiveButton(R.string.leave_room, new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                try {
                    mService.sendMessage(JSONParser.leaveRoomRequest(roomID, username));
                } catch (JSONException e) {
                    e.printStackTrace();
                }

                // TODO: Handle this when response received
                roomDB.deleteRoom(roomID);
                roomList.remove(position);
                roomAdapter.notifyDataSetChanged();
            }
        });
        builder.setNegativeButton("Cancel", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                // Cancel
            }
        });
        builder.create();
        builder.show();
    }

    public void onAddRoomClick(final int createOrJoin) {
        actionMenu.collapse();
        LayoutInflater inflater = LayoutInflater.from(context);
        View dialogView = inflater.inflate(R.layout.dialog_add_room, null);

        AlertDialog.Builder builder = new AlertDialog.Builder(context);
        builder.setView(dialogView);

        final EditText inputRoom = (EditText) dialogView.findViewById(R.id.input_room);
        if (getResources().getString(createOrJoin).equals(getResources().getString(R.string.join_room))) {
            inputRoom.setHint(R.string.enter_room_id);
            inputRoom.setInputType(InputType.TYPE_CLASS_NUMBER);
        }

        builder.setPositiveButton(createOrJoin, new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                String nameOrId = inputRoom.getText().toString();

                if (getResources().getString(createOrJoin).equals(getResources().getString(R.string.join_room))) {
                    joinRoomRequest(Integer.parseInt(nameOrId));
                } else {
                    createRoomRequest(nameOrId);
                }
            }
        });

        AlertDialog alertDialog = builder.create();
        alertDialog.getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_VISIBLE);
        alertDialog.show();
    }

    private void createRoomRequest(String roomName) {
        try {
            mService.sendMessage(JSONParser.createRoomRequest(roomName, username));
        } catch (JSONException e) {
            e.printStackTrace();
        }
    }

    private void joinRoomRequest(int roomid) {
        try {
            mService.sendMessage(JSONParser.joinRoomRequest(roomid, username));
        } catch (JSONException e) {
            e.printStackTrace();
        }
    }

    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.main_menu, menu);
        return true;
    }

    public void logout() {
        SharedPreferences loginInfo = getSharedPreferences(LOGIN_PREFS, 0);
        SharedPreferences.Editor loginEdit = loginInfo.edit();
        loginEdit.putString("username", "");
        loginEdit.putBoolean("loggedIn", false);
        loginEdit.commit();
    }

    public void onLogoutClick() {

        AlertDialog.Builder builder = new AlertDialog.Builder(context);
        builder.setMessage(R.string.logout_question);
        builder.setPositiveButton("Yes", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                // TODO: Send logout request
                try {
                    mService.sendMessage(JSONParser.signoutRequest(username));
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        });
        builder.setNegativeButton("No", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                // Cancel
            }
        });
        builder.create();
        builder.show();
    }

    public boolean onOptionsItemSelected(MenuItem item) {
        int id = item.getItemId();

        if (id == R.id.action_logout) {
            onLogoutClick();
            return true;
        }
        return super.onOptionsItemSelected(item);
    }
}
