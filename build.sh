#!/bin/bash
go mod tidy && CGO_ENABLED=0 go build -o devbox . && sudo mv devbox  /usr/local/bin/ && sudo devbox

