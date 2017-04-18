package com.myplace.myplace;

import android.content.Intent;
import android.os.Bundle;
import android.support.annotation.NonNull;
import android.support.annotation.Nullable;
import android.support.design.widget.FloatingActionButton;
import android.support.design.widget.Snackbar;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.ListView;
import android.widget.TextView;

import java.util.ArrayList;

public class MainActivity extends AppCompatActivity {

    Toolbar toolbar;
    ListView listView;
    ArrayList<Room> roomList = new ArrayList<>();

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        Toolbar toolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);
        //noinspection ConstantConditions
        getSupportActionBar().setTitle(getResources().getString(R.string.app_name));

        roomList.add(new Room("Ett rum"));
        roomList.add(new Room("två rum"));

        listView = (ListView) findViewById(R.id.roomList);


        ArrayAdapter<Room> adapter = new ArrayAdapter<Room>(MainActivity.this, android.R.layout.simple_list_item_2, android.R.id.text1, roomList) {

            @NonNull
            @Override
            public View getView(int position, @Nullable View convertView, @NonNull ViewGroup parent) {
                View view = super.getView(position, convertView, parent);
                TextView text1 = (TextView) view.findViewById(android.R.id.text1);
                TextView text2 = (TextView) view.findViewById(android.R.id.text2);

                text1.setText(roomList.get(position).getName());
                text2.setText(roomList.get(position).getLastMessage());
                return view;
            }
        };

        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                Intent intent = new Intent(MainActivity.this, MessageActivity.class);
                intent.putExtra("RoomName", roomList.get(position).getName());
                //intent.putExtra("SongName", listView.getItemAtPosition(position).toString());
                startActivity(intent);
            }
        });
/*

        RoomAdapter<Room> adapter = new RoomAdapter<Room>(MainActivity.this,
                android.R.layout.simple_list_item_2,
                new ArrayList<Room>(roomList));
*/

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
