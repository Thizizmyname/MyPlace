package com.myplace.myplace;

import android.content.Context;
import android.content.SharedPreferences;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import com.myplace.myplace.models.Message;

import java.util.ArrayList;

import static com.myplace.myplace.LoginActivity.LOGIN_PREFS;

public class MessageAdapter extends ArrayAdapter<Message> {
    private static final int TYPE_INCOMING = 0;
    private static final int TYPE_OUTGOING = 1;
    private static final int TYPE_MAX_COUNT = TYPE_OUTGOING + 1;

    ArrayList<Message> messages;

    public MessageAdapter(Context context, ArrayList<Message> messages) {
        super(context, 0, messages);
        this.messages = messages;
    }

    @Override
    public int getItemViewType(int position) {
        Message message = getItem(position);

        String username = "N/A";
        SharedPreferences loginInfo = getContext().getSharedPreferences(LOGIN_PREFS, 0);
        if(message.getName().equals(loginInfo.getString("username", username))) { // TO TEST INCOMING AND OUTGOING
            return TYPE_OUTGOING;
        } else {
            return TYPE_INCOMING;
        }
    }

    @Override
    public int getViewTypeCount() {
        return TYPE_MAX_COUNT;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        Message message = getItem(position);

        int itemType = getItemViewType(position);
        if (convertView == null) {
            if (itemType == TYPE_INCOMING) {
                convertView = LayoutInflater.from(getContext()).inflate(R.layout.message_bubble_left, parent, false);
            } else {
                convertView = LayoutInflater.from(getContext()).inflate(R.layout.message_bubble_right, parent, false);
            }
        }

        if (itemType == TYPE_INCOMING) {
            TextView name = (TextView) convertView.findViewById(R.id.textName);
            name.setText(message.getName());
        }

        TextView text = (TextView) convertView.findViewById(R.id.textMessage);
        TextView date = (TextView) convertView.findViewById(R.id.textDate);

        date.setText(message.getShortTime());
        text.setText(message.getText());

        return convertView;
    }

    public void updateData(ArrayList<Message> list) {
        this.messages.clear();
        messages.addAll(list);
        this.notifyDataSetChanged();
    }
}
