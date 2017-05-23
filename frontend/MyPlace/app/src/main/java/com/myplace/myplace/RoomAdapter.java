package com.myplace.myplace;

import android.content.SharedPreferences;
import android.graphics.Typeface;
import android.support.annotation.NonNull;
import android.support.annotation.Nullable;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import com.myplace.myplace.models.Room;
import com.myplace.myplace.models.RoomInfo;

import java.util.ArrayList;

import static com.myplace.myplace.LoginActivity.LOGIN_PREFS;

/**
 * Created by alexis on 2017-04-18.
 */

class RoomAdapter extends ArrayAdapter<RoomInfo> {

    ArrayList<RoomInfo> rooms;

    public RoomAdapter(MainActivity context, ArrayList<RoomInfo> rooms) {
        super(context, 0, rooms);
        this.rooms = rooms;

    }

    @NonNull
    @Override
    public View getView(int position, @Nullable View convertView, @NonNull ViewGroup parent) {
        RoomInfo roomInfo = getItem(position);

        if (convertView == null) {
            convertView = LayoutInflater.from(getContext()).inflate(R.layout.custom_room_list_item, parent, false);
        }

        TextView text1 =  (TextView) convertView.findViewById(R.id.r_title);
        TextView text2 =  (TextView) convertView.findViewById(R.id.r_name);
        TextView text3 =  (TextView) convertView.findViewById(R.id.r_subtitle);

        if (roomInfo != null) {
            if (roomInfo.getLastMessage() != null && roomInfo.getLastMessage().getId() != roomInfo.getLastMsgRead()) {
                text1.setTypeface(Typeface.DEFAULT_BOLD);
                text2.setTypeface(Typeface.DEFAULT_BOLD);
                text3.setTypeface(Typeface.DEFAULT_BOLD);
            } else {
                text1.setTypeface(Typeface.DEFAULT);
                text2.setTypeface(Typeface.DEFAULT);
                text3.setTypeface(Typeface.DEFAULT);
            }

            text1.setText(roomInfo.getName());
            String sender = roomInfo.getLastSender();

            String username = "N/A";
            SharedPreferences loginInfo = getContext().getSharedPreferences(LOGIN_PREFS, 0);
            if (sender.equals(loginInfo.getString("username", username))) {
                String name = getContext().getString(R.string.sender_you) + ": ";
                text2.setText(name);
            } else if (!sender.equals("")) {
                String name = sender+": ";
                text2.setText(name);
            } else {
                text2.setText(sender);
            }
            text3.setText(roomInfo.getLastMessageText());
        }

        return convertView;
    }

    public void updateData(ArrayList<RoomInfo> list) {
        this.rooms.clear();
        rooms.addAll(list);
    }
}
