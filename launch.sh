#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
sudo docker run -dit --rm -v $DIR/scripts:/scripts -v $DIR/usage:/usage --cap-add=NET_ADMIN ns
