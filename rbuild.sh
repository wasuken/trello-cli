#!/bin/bash

export GOOS=linux
export GOARCH=arm
export GOARM=7
go build
