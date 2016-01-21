## Concurrent File Server
===========================

First assignment for course CS733, Spring-2015  
Submitted by: Mihir J. Vegad, Roll: 143050073

### Introduction
=================
* This is a single server for a versioned files to be stored and retrieved by the server   
* The server listens on port 8080.  
* It can handle multiple clients concurrently.  
* Automated testing is provided with Go's testing framework

### Installation Instructions
==============================
<code>go get </code> github.com/MJVegad/cs733/assignment1

Two types of files are supposed to be there <br/>
1. mysingleserver.go contains the code where all the commands are implemented and where server listens for the request <br/>
2. All _test.go contains all the test cases including commands which are fired concurrently evaluating all the necessary scenarios

To run the program only below command is needed (assuming the current directory is set to the assignment1 which has the go files)
<br/><code>go test -race</code>


### How to use?
================
Run the server by "go run mysingleserver.go" in one terminal from the assignment1 directory.
Run "telnet localhost 8080" on a different terminal to connect with the server as a telnet client.
The server supports multiple client connections by executing same command on different terminals and allows concurrent execution of these clients.
Once a connection is established one of the commands given in specification below can be run by a user through a client.

### Specification
==================
* Write: create a file, or update the file’s contents if it already exists.
```
write <filename> <numbytes> [<exptime>]\r\n
 <content bytes>\r\n
```
exptime field is optional, it signifies the time in seconds after which the file on server 
is suppose to expire. If left unspecified or if its value is set to zero then file does not expires.

The server responds with the following:

```

OK <version>\r\n

``````
where version is a unique 64‐bit number (in decimal format) assosciated with the
filename.

* Read: Given a filename, retrieve the corresponding file:
```
read <filename>\r\n
```
The server responds with the following format if file is present at the server.
```
CONTENTS <version> <numbytes> <exptime> \r\n
 <content bytes>\r\n  
```
Here ```<exptime> ```is the remaining time in seconds left for the file after which it will expire. Zero v alue indicates that the file won't expire.

If the file is not present at the server or the file has expired then the server response is:
```
ERR_FILE_NOT_FOUND\r\n
```

* Compare and swap (cas): This replaces the old file contents with the new content
provided the version is still the same.
```
cas <filename> <version> <numbytes> [<exptime>]\r\n
 <content bytes>\r\n
```
exptime is optional and means the same thing as in "write" command.

The server responds with the new version if successful :-
```
OK <version>\r\n
```
If the file isn't found on the server or the file has expired then server response is:-
```
ERR_FILE_NOT_FOUND\r\n
```
If the version provided in the "cas" command does not match the version of the file, server
response is:-
```
ERR_VERSION <version>\r\n
```

* Delete file
```
delete <filename>\r\n
```
Server response (if successful)
```
OK\r\n
```

If the file isn't found on the server or the file has expired then server response is:-
```
ERR_FILE_NOT_FOUND\r\n
```

Apart from above mentioned errors if the above commands are not specified in proper format as mentioned, then server response is:-
```
ERR_CMD_ERR\r\n
```
After sending the above response to client, server terminates the client connection and thus forcing the client to connect again.


### Testing
============
For running the automated testing procedure, ensure that the server is not running and any other process is not busy at port 8080.
Run "go test -race" to run the tests. It would take some time.

For a single/multiple client server communication most of the scenarios along with corner cases are tested. This program runs serially.

For each of the test cases the client sends a command to the server, reads the response and checks the result with the specified expected result mentioned in the test case.


### Programming Details
========================
* A structure is is used to store values corresponding to the file into a map. Key of map is filename and value is the structure.
* Read/Write mutex are used to handle the race conditions created by the concurrent access of clients on the shared data structure.
* The server accepts connection on port 8080 and for every client handles its connection on a different thread.

### Todo
1. To make the file server storage persistent suing leveldb

(source: github.com/jain7aman/cs733/assignment1)