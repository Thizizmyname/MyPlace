package com.myplace.myplace;

import android.content.Context;
import android.text.Layout;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import java.util.ArrayList;

/**
 * Created by jesper on 2017-04-06.
 */

public class MessageAdapter extends ArrayAdapter<Message> {
    public MessageAdapter(Context context, ArrayList<Message> messages) {
        super(context, 0, messages);
    }

    public View getView(int position, View convertView, ViewGroup parent){
        Message message = getItem(position);

        TextView name = (TextView) convertView.findViewById(R.id.textName);
        TextView text = (TextView) convertView.findViewById(R.id.textMessage);

        if(convertView == null) {
            int layout;
            if (message.name == "Anders") {
                layout = R.layout.message_bubble_right;
            } else {
                layout = R.layout.message_bubble_left;
            }
            convertView = LayoutInflater.from(getContext()).inflate(layout, parent, false);
        }

        name.setText(message.name);
        text.setText(message.text);

        return convertView;
    }
}
