#!/bin/bash
SRC="/gopath/src/github.com/CovenantSQL/GNTE"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
docker run --rm -it -v $DIR:$SRC gnte $SRC/scripts/gobuild.sh
