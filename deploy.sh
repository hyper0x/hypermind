#!/bin/sh

SERVER_PORT="9091"
PID=`ps -ef | grep "hypermind" | grep -v "cronolog" | grep -v "grep" | awk '{print $2}'`
if [ -z $PID ]; then
	PID=`lsof -i:$SERVER_PORT | grep "a.out" | awk '{print $2}'`
fi
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
echo "Building the program ..."
go build
$DIR/hypermind -port=$SERVER_PORT 2>&1 | cronolog "$LOG_DIR/hypermind.log.%Y-%m-%d" &
echo "The server has been launched. (port=$SERVER_PORT)"
