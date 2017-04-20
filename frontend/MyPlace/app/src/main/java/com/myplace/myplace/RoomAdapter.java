package com.myplace.myplace;

import android.support.annotation.NonNull;
import android.support.annotation.Nullable;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-04-18.
 */

class RoomAdapter<T> extends ArrayAdapter<Room> {

    ArrayList<Room> rooms;

    public RoomAdapter(MainActivity mainActivity, int simple_list_item_2, ArrayList<Room> rooms) {
        super(mainActivity, simple_list_item_2);
        this.rooms = rooms;

    }

    @NonNull
    @Override
    public View getView(int position, @Nullable View convertView, @NonNull ViewGroup parent) {

        View view = super.getView(position, convertView, parent);
        TextView text1 =  (TextView) view.findViewById(android.R.id.text1);
        TextView text2 =  (TextView) view.findViewById(android.R.id.text2);

        text1.setText(rooms.get(position).getName());
        text2.setText(rooms.get(position).getLastMessage());


        return super.getView(position, convertView, parent);
    }
}
