#!/bin/sh

PID=`ps -ef | grep "go-web-demo" | grep -v grep|awk '{print $2}'`
kill $PID
echo "The go-web-demo has stopped!"
./go-web-demo 2>&1 | cronolog "go-web-demo.log.%Y-%m-%d" &
echo "The go-web-demo has been launched."
