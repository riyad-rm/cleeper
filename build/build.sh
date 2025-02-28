#!/bin/bash

if [ ! -d "bin" ]; then
  mkdir bin
fi

if [ -f bin/bootstrap ]; then
    rm bin/bootstrap
fi

if [ -f bin/cleeper.zip ]; then
    rm bin/cleeper.zip
fi

cd ../src/
go mod tidy
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../build/bin/bootstrap
cd ../build/bin
strip -s bootstrap
zip cleeper.zip bootstrap