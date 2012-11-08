#!/bin/sh

PID=`ps -ef | grep "go-web-demo" | grep "server" | grep -v grep | awk '{print $2}'`
kill $PID
echo "The go-web-demo has stopped!"
go build server
chmod +rx server
DIR=`pwd`
$DIR/server 2>&1 | cronolog "go-web-demo.log.%Y-%m-%d" &
echo "The go-web-demo has been launched."
