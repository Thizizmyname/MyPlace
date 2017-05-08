package com.myplace.myplace;

import android.support.annotation.NonNull;
import android.support.annotation.Nullable;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import com.myplace.myplace.models.Room;

import java.util.ArrayList;

/**
 * Created by alexis on 2017-04-18.
 */

class RoomAdapter extends ArrayAdapter<Room> {

    ArrayList<Room> rooms;

    public RoomAdapter(MainActivity context, ArrayList<Room> rooms) {
        super(context, 0, rooms);
        this.rooms = rooms;

    }

    @NonNull
    @Override
    public View getView(int position, @Nullable View convertView, @NonNull ViewGroup parent) {
        Room room = getItem(position);


        if (convertView == null) {
            convertView = LayoutInflater.from(getContext()).inflate(R.layout.custom_room_list_item, parent, false);
        }

        TextView text1 =  (TextView) convertView.findViewById(R.id.r_title);
        TextView text2 =  (TextView) convertView.findViewById(R.id.r_subtitle);

        if (room != null) {
            text1.setText(room.getName());
            text2.setText(room.getLastMessage(getContext()));
        }





        return convertView;
    }
}
