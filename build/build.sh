#!/bin/bash

rm bin/bootstrap
rm bin/cleeper.zip
cd ../src/
go mod tidy
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../build/bin/bootstrap
cd ../build/bin
strip -s bootstrap
zip cleeper.zip bootstrap