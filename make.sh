#!/usr/bin/env bash
#gofmt -tabwidth 4 -w ./
rm *.8 2>/dev/null
8g util.go
8g client.go
8g server.go
8l -o p2pclient client.8 util.8
8l -o p2pserver server.8 util.8
