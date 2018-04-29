#!/bin/sh

if [ $USER = 'root' ];
then
    docker build -t ns .
else
    sudo docker build -t ns .
fi
