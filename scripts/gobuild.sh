#!/bin/sh
export GOPATH=/gopath
cd /gopath/src/github.com/CovenantSQL/GNTE
go run main.go $*
