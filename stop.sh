#!/bin/bash
server="myServer"
pids=`ps -ef |grep ./$server |grep -v "grep"   |awk '{print $2}'`
echo "the pid list:" $pids
kill -9 $pids
echo "finish stop the server cluster..."

