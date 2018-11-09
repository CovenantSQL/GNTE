#!/bin/bash

param=$1
filter=$2

generate() {

    SRC="/gopath/src/github.com/CovenantSQL/GNTE"
    export DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    CLEAN=$DIR/scripts/clean.sh

    if [ -f $CLEAN ];then
        $CLEAN
        rm -rf $CLEAN
    fi

    docker run --rm -it -v $DIR:$SRC gnte $SRC/scripts/gobuild.sh $param

    $DIR/scripts/launch.sh
}

get_containers() {
    if [ -n $filter ]; then
        containers="$(docker ps --format '{{.Names}}' --filter 'network=CovenantSQL_testnet' --filter name=$filter)"
    else
        containers="$(docker ps --format '{{.Names}}' --filter 'network=CovenantSQL_testnet')"
    fi
    echo $containers
}

stopone() {
    containers=`get_containers`
    for i in $containers; do
        array=("${array[@]}" $i)
    done
    len=${#array[@]}
    if [ 0 -eq $len ]; then
        return
    fi
    num=$(date +%s)
    ((rand=num%len))
    echo "Stopping ${array[$rand]}"
    docker stop ${array[$rand]}
}

stopall() {
    containers=`get_containers`
    for i in $containers; do
        echo "Stopping $i"
        docker stop $i
    done
}

startall() {
    containers="$(docker ps --format '{{.Names}}' --filter 'network=CovenantSQL_testnet' --filter status=exited)"

    for i in $containers; do
        echo "Starting $i"
        docker start $i
    done
}

case $param in
    "stopone")
        stopone
        ;;
    'stopall')
        stopall
        ;;
    'startall')
        startall
        ;;
    *)
        echo "Generate GNTE and running"
        generate
        ;;
esac
