package com.myplace.myplace;

import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.WindowManager;
import android.widget.AdapterView;
import android.widget.EditText;
import android.widget.ListView;
import com.myplace.myplace.models.Room;
import java.util.ArrayList;
import android.view.WindowManager;
import android.widget.EditText;
import android.widget.Toast;

import com.getbase.floatingactionbutton.FloatingActionButton;
import com.getbase.floatingactionbutton.FloatingActionsMenu;



import static com.myplace.myplace.R.id.action_create;
import static com.myplace.myplace.R.id.action_join;

public class MainActivity extends AppCompatActivity {
    private TCPClient mTcpClient;
    final Context context = this;
    Toolbar toolbar;

    FloatingActionsMenu actionMenu;
    ListView listView;
    ArrayList<Room> roomList = null;
    public static RoomAdapter roomAdapter = null;

    //Defines the database
    public static RoomDbHelper roomDB = null;

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

        roomAdapter = new RoomAdapter(MainActivity.this, roomList);
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

    public class ConnectTask extends AsyncTask<String,String,TCPClient> {

        @Override
        protected TCPClient doInBackground(String... message) {

            //we create a TCPClient object and
            mTcpClient = new TCPClient(new TCPClient.OnMessageReceived() {
                @Override
                //here the messageReceived method is implemented
                public void messageReceived(String message) {
                    //this method calls the onProgressUpdate
                    publishProgress(message);
                    Log.d("Message", message);
                }
            });
            mTcpClient.run();

            return null;
        }

        @Override
        protected void onProgressUpdate(String... values) {
            super.onProgressUpdate(values);
            Log.d("values", values[0]);

            // TODO: Handle response from server
        }
    }

    public void onThreadClick(int position) {
        Intent intent = new Intent(MainActivity.this, MessageActivity.class);
        intent.putExtra("Room", roomList.get(position));
        startActivity(intent);
        overridePendingTransition(R.anim.push_left_in, R.anim.push_left_out);
    }

    public void onThreadLongClick(final int position) {
        final String roomName = roomList.get(position).getName();

        AlertDialog.Builder builder = new AlertDialog.Builder(context);
        builder.setMessage("Do you want to leave "+roomName+"?");
        builder.setPositiveButton(R.string.leave_room, new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                roomDB.deleteRoom(roomName);
                roomList.remove(position);
                roomAdapter.notifyDataSetChanged();

                //TODO: Change below code to JSON-request
                TCPClient.request = "Leave room "+roomName;
                new ConnectTask().execute("");
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

        builder.setPositiveButton(createOrJoin, new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                String roomName = inputRoom.getText().toString();

                //TODO: Change below string to JSON-request

                roomDB.createRoomTable(roomName);
                roomList.add(new Room(0,roomName));
                roomAdapter.notifyDataSetChanged();

                TCPClient.request = getResources().getString(createOrJoin)+roomName;
                new ConnectTask().execute("");
            }
        });

        AlertDialog alertDialog = builder.create();
        alertDialog.getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_VISIBLE);
        alertDialog.show();
    }
}
