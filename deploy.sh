#!/bin/sh

PID=`ps -ef | grep "hypermind" | grep "server" | grep -v grep | awk '{print $2}'`
if [ -z $PID ]; then
	echo "The server has yet to launch."
else
	kill $PID
	echo "The server has stopped!"
fi
DIR=`pwd`
LOG_DIR=$DIR/logs
if [ ! -d "$LOG_DIR" ]; then
	echo "Creating dir '$LOG_DIR' ..."
	mkdir "$LOG_DIR"
fi
go run server.go 2>&1 | cronolog "$LOG_DIR/hypermind.log.%Y-%m-%d" &
echo "The server has been launched."
