# gofinger

This repo holds a Finger-like daemon and client. The server is pretty simple: you send it a TCP packet containing a username, it looks up that username on the system, and then prints back information about the user. The Internet really _was_ a simpler place long ago.

## Compiling

To build the server and client, run `make`. To remove artifacts, run `make clean`.

## Running the Server

```
$ cd fingerd
$ ./fingerd
2019/06/05 00:41:48 listening on 127.0.0.1:4500
```

## Running the Client

```
$ cd finger
$ ./finger root

username: root  name: System Administrator
directory: /var/root
No plan file.
```

## Notes

* This _should_ have enough error handling and defer-style closing to be safe. Please let me know if it doesn't!
* The server will only bind to the loopback interface. That's intentional--I didn't want this program to be accessible outside this machine.
