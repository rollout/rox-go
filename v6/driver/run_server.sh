#!/bin/bash -x
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd -P "$DIR"
DIR=`pwd`
export PORT=$1
go get
go build
nohup ./driver 1>"$DIR"/log_"$PORT".out 2>&1 &
while true ; do
  curl -p http://127.0.0.1:$PORT/status-check && exit
  sleep 1
done
