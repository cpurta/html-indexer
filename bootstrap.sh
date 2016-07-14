#!/bin/bash

# ensure that Go is installed if not print error message
go version 1>/dev/null
if [ $? -ne 0 ]; then
    echo "Go does not appear to be installed. Please install and try again."
    exit 1
fi

# go get all of our go dependencies
echo "Getting project dependencies..."
go get github.com/go-sql-driver/mysql
go get github.com/influxdata/influxdb/client/v2
go get gopkg.in/redis.v2

echo "Building hmtl-indexer..."
go build -o html-indexer .

exit 0
