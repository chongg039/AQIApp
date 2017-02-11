# Makefile for AQI app
#
# Author: coldriver
#
# Targets:
#	build: Builds the code
SERVER=server
CRAWLER=crawler

default: build

build: 
	go build -o ../build/${SERVER} ./server/*.go
	go build -o ../build/${CRAWLER} ./crawler/*.go
