#!/bin/bash

#build and run
killall chat_svc
go build -o chat_svc main.go
nohup ./chat_svc > /dev/null 2>&1 &

