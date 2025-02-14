#!/bin/bash

cd bin
go build -ldflags="-s -w" -o cleeper ../../src
strip -s cleeper
zip cleeper.zip cleeper
