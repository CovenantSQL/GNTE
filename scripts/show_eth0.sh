#!/bin/sh

if [ -z $1 ]
then
    tc -s qdisc show dev eth0
elif [ $1 = '-w' ]
then
    watch -n 1 tc -s qdisc show dev eth0
fi
