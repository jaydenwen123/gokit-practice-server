#!/bin/bash
echo "begin to complie the program..."
server="myServer"
go build -o $server
echo "finish the complie...."
echo "begin to start the server cluster...."
./$server --name "userService" --port 2345 &
./$server --name "userService" --port 2346 &
./$server --name "userService" --port 2347 &
./$server --name "userService" --port 2348 &
./$server --name "userService" --port 2349 &
./$server --name "userService" --port 2350 &
echo "the server cluster finish started..."
