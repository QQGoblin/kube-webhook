#!/bin/bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOPROXY=https://goproxy.cn,direct go build -o ./output/kube-webhook  ./main.go
