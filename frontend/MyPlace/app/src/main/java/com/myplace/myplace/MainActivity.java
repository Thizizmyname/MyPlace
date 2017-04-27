package com.myplace.myplace;

import android.content.DialogInterface;
import android.content.Intent;
import android.database.sqlite.SQLiteDatabase;
import android.os.Bundle;
import android.support.annotation.NonNull;
import android.support.annotation.Nullable;
import android.support.design.widget.Snackbar;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.TextView;
import android.graphics.drawable.ShapeDrawable;
import android.graphics.drawable.shapes.OvalShape;
import com.getbase.floatingactionbutton.FloatingActionButton;
import com.getbase.floatingactionbutton.FloatingActionsMenu;


import java.util.ArrayList;

import static android.R.id.input;
import static com.myplace.myplace.R.id.action_create;
import static com.myplace.myplace.R.id.input_room;

public class MainActivity extends AppCompatActivity {

    Toolbar toolbar;
    ListView listView;
    ArrayList<Room> roomList = null;
    public static RoomAdapter adapter = null;

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

        Room r1 = new Room("Rum 1");
        roomDB.createRoomTable(r1.getName());


        Room r2 = new Room("Rum 2");
        roomDB.createRoomTable(r2.getName());


        roomList = roomDB.getRoomList();

        listView = (ListView) findViewById(R.id.roomList);


        adapter = new RoomAdapter(MainActivity.this, roomList);


        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                Intent intent = new Intent(MainActivity.this, MessageActivity.class);
                intent.putExtra("Room", roomList.get(position));
                startActivity(intent);
            }
        });


        listView.setAdapter(adapter);

        final View actionCreate = findViewById(action_create);

        actionCreate.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AlertDialog.Builder builder = new AlertDialog.Builder(getApplicationContext());
                builder.setView(R.layout.dialog_add_room);

                builder.setPositiveButton("Create room", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        EditText input = (EditText) findViewById(R.id.input_room);
                        String roomName = input.getText().toString();

                        roomDB.createRoomTable(roomName);
                    }
                });
                builder.show();
            }
        });
    }


}
