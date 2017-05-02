package com.myplace.myplace;

import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.AdapterView;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.Toast;

import com.getbase.floatingactionbutton.FloatingActionButton;
import com.getbase.floatingactionbutton.FloatingActionsMenu;


import java.util.ArrayList;

import static com.myplace.myplace.R.id.action_create;
import static com.myplace.myplace.R.id.action_join;

public class MainActivity extends AppCompatActivity {

    final Context context = this;
    Toolbar toolbar;
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

        //TODO Tests for viewing temporary items

        roomList = roomDB.getRoomList();

        listView = (ListView) findViewById(R.id.roomList);


        roomAdapter = new RoomAdapter(MainActivity.this, roomList);
        listView.setAdapter(roomAdapter);

        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                Intent intent = new Intent(MainActivity.this, MessageActivity.class);
                intent.putExtra("Room", roomList.get(position));
                startActivity(intent);
                overridePendingTransition(R.anim.push_left_in, R.anim.push_left_out);
            }
        });

        listView.setOnItemLongClickListener(new AdapterView.OnItemLongClickListener() {
            @Override
            public boolean onItemLongClick(AdapterView<?> parent, View view, final int position, long id) {
                final String roomName = roomList.get(position).getName();

                AlertDialog.Builder builder = new AlertDialog.Builder(context);
                builder.setMessage("Do you want to leave "+roomName+"?");
                builder.setPositiveButton("Leave room", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        roomDB.deleteRoom(roomName);
                        roomList.remove(position);
                        roomAdapter.notifyDataSetChanged();
                        //TODO: Send leave room request
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
                return true;
            }
        });


        final FloatingActionsMenu actionMenu = (FloatingActionsMenu) findViewById(R.id.action_menu);
        //OnClick for createRoom
        final FloatingActionButton actionCreate = (FloatingActionButton) findViewById(action_create);

        actionCreate.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                actionMenu.collapse();
                LayoutInflater inflater = LayoutInflater.from(context);
                View dialogView = inflater.inflate(R.layout.dialog_add_room, null);

                AlertDialog.Builder builder = new AlertDialog.Builder(context);
                builder.setView(dialogView);

                final EditText inputRoom = (EditText) dialogView.findViewById(R.id.input_room);

                builder.setPositiveButton("Create room", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        String roomName = inputRoom.getText().toString();

                        //TODO: Send request to create a new room
                        Room room = new Room(roomName);
                        roomAdapter.add(room);
                        roomDB.createRoomTable(roomName);
                    }
                });

                AlertDialog alertDialog = builder.create();
                alertDialog.show();
            }
        });


        //OnClick for joinRoom
        final FloatingActionButton actionJoin = (FloatingActionButton) findViewById(action_join);

        actionJoin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                actionMenu.collapse();
                LayoutInflater inflater = LayoutInflater.from(context);
                View dialogView = inflater.inflate(R.layout.dialog_add_room, null);

                AlertDialog.Builder builder = new AlertDialog.Builder(context);
                builder.setView(dialogView);

                final EditText inputRoom = (EditText) dialogView.findViewById(R.id.input_room);

                builder.setPositiveButton("Join room", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        String roomName = inputRoom.getText().toString();

                        Toast toast = Toast.makeText(MainActivity.this, "Not yet implemented", Toast.LENGTH_SHORT);
                        toast.show();

                        //TODO: Send request to join an existing room

                    }
                });

                AlertDialog alertDialog = builder.create();
                alertDialog.show();
            }
        });
    }


}
