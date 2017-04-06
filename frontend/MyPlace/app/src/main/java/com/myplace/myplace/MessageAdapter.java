package com.myplace.myplace;

import android.content.Context;
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

        if(convertView == null) {
            convertView = LayoutInflater.from(getContext()).inflate(R.layout.message_bubble, parent, false);
        }

        TextView name = (TextView) convertView.findViewById(R.id.textName);
        TextView text = (TextView) convertView.findViewById(R.id.textMessage);

        name.setText(message.name);
        text.setText(message.text);

        return convertView;
    }
}
