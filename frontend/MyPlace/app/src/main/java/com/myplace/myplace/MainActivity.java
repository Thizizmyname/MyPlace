package com.myplace.myplace;

import android.content.Intent;
import android.os.Bundle;
import android.support.design.widget.FloatingActionButton;
import android.support.design.widget.Snackbar;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ListView;

import com.myplace.myplace.models.Room;

import java.util.ArrayList;

import org.json.JSONObject;

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
                intent.putExtra("RoomName", roomList.get(position).getName());
                startActivity(intent);
            }
        });


        listView.setAdapter(adapter);

        FloatingActionButton fab = (FloatingActionButton) findViewById(R.id.fab);
        fab.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Snackbar.make(view, "Replace with your own action", Snackbar.LENGTH_LONG)
                        .setAction("Action", null).show();
            }
        });
    }

}