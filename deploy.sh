#!/bin/bash

# mac
go build -o mac main.go

# 64位win
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -o win.exe main.go
