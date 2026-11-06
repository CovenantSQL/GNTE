#!/bin/sh
cd /go/src/github.com/CovenantSQL/GNTE
go run -mod=vendor main.go $*
