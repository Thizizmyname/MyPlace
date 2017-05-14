package com.myplace.myplace;

import android.app.IntentService;
import android.app.Service;
import android.content.Intent;
import android.os.AsyncTask;
import android.os.Binder;
import android.os.IBinder;
import android.support.annotation.Nullable;
import android.util.Log;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import java.net.InetAddress;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.Random;

/**
 * Created by alexis on 2017-05-10.
 */

public class TCPService extends Service {
    // Binder given to clients
    private final IBinder mBinder = new TCPBinder();
    // Random number generator
    private final Random mGenerator = new Random();

    private String serverMessage;

    public static final String SERVERIP = "10.0.2.2";
    public static final int SERVERPORT = 1337;
    private static Socket socket;
    private static PrintWriter out;
    BufferedReader in;

    private OnMessageReceived mMessageListener = null;
    private boolean mRun = false;
    private boolean pauseListener = false;


    /**
     * Class used for the client Binder.  Because we know this service always
     * runs in the same process as its clients, we don't need to deal with IPC.
     */
    public class TCPBinder extends Binder {
        TCPService getService() {
            // Return this instance of LocalService so clients can call public methods
            Log.e("TCP Service", "Started Service");
            return TCPService.this;
        }
    }

    @Override
    public IBinder onBind(Intent intent) {
        Log.e("TCP Service", "In onBind");
        return mBinder;
    }

    /** method for clients */
    public int giveRandomNumber() {
        return mGenerator.nextInt(100);
    }

    public void sendMessage(final String message){

        AsyncTask<Void, Void, Boolean> sendThread = new AsyncTask<Void, Void, Boolean>() {
            @Override
            protected Boolean doInBackground(Void... params) {

                //try {
                //out = new PrintWriter(new BufferedWriter(new OutputStreamWriter(socket.getOutputStream())), true);
                pauseListener = true;

                Log.e("TCP Service", "Sending: " + message);
                if (out != null && !out.checkError()) {
                    out.println(message);
                    Log.d("TCP Client", "Message: " + message);
                    out.flush();
                    //out.close();
                }

                pauseListener = false;
//                } catch (IOException e) {
//                    e.printStackTrace();
//                }



                return null;
            }
        };
        sendThread.executeOnExecutor(AsyncTask.THREAD_POOL_EXECUTOR);


    }

    private void runListener() {
        while (mRun) {
            //Log.e("TCP Client", "C: I got to the while loop!");
            if (!pauseListener) {
                try {
                    serverMessage = in.readLine();


                    Log.e("TCP Service", "C: serverMessage = " + serverMessage);

                    if (serverMessage != null && mMessageListener != null) {
                        Log.e("TCP Client", "C: serverMessage = " + serverMessage);
                        //call the method messageReceived from MyActivity class
                        mMessageListener.messageReceived(serverMessage);
                    } else {
                        serverMessage = null;
                    }
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }

        Log.e("TCP Client", "C: run = " + mRun);

        Log.e("RESPONSE FROM SERVER", "S: Received Message: '" + serverMessage + "'");

    }

    public void setUpConnection() {

        AsyncTask<Void, Void, Boolean> connectionThread = new AsyncTask<Void, Void, Boolean>() {

            @Override
            protected Boolean doInBackground(Void... arg0) {

                mRun = true;

                try {
                    //here you must put your computer's IP address.
                    InetAddress serverAddr = InetAddress.getByName(SERVERIP);

                    Log.e("TCP Client", "C: Connecting...");

                    //create a socket to make the connection with the server
                    socket = new Socket(serverAddr, SERVERPORT);
                    //out = new PrintWriter(new BufferedWriter(new OutputStreamWriter(socket.getOutputStream())), true);

                    try {

                        out = new PrintWriter(new BufferedWriter(new OutputStreamWriter(socket.getOutputStream())), true);


                        //receive the message which the server sends back
                        in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
                        Log.e("TCP Client", "C: received = " + in);
                        Log.e("TCP Client", "C: run = " + mRun);
                        //in this while the client listens for the messages sent by the server

                        runListener();

                    } catch(Exception e){
                        Log.e("TCP", "S: Error", e);
                    } finally{
                        //the socket must be closed. It is not possible to reconnect to this socket
                        // after it is closed, which means a new socket instance has to be created.
                        socket.close();
                        Log.d("TCP Client", "Socket closed.");
                        serverMessage = null;
                    }


                } catch (UnknownHostException e) {
                    e.printStackTrace();
                } catch (IOException e) {
                    e.printStackTrace();
                }

                return null;
            }
        };
        connectionThread.executeOnExecutor(AsyncTask.THREAD_POOL_EXECUTOR);
    }

    //Declare the interface. The method messageReceived(String message) will must be implemented in the MyActivity
    //class at on asynckTask doInBackground
    public interface OnMessageReceived {
        void messageReceived(String message);
    }
}
