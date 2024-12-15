#!/bin/bash

COUNT_TO="$1"

i=1
while [ $i -le $COUNT_TO ]; do
    echo "counting: $i"
    i=$((i+1))
    sleep 1
done
