#!/bin/bash

SLEEP="$1"
COUNT_TO="$2"

i=1
while [ $i -le $COUNT_TO ]; do
    echo "counting: $i"
    i=$((i+1))
    sleep $SLEEP
done
