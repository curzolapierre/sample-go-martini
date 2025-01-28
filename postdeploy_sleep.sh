#!/bin/bash

# This script is used to sleep 25 minutes and echo each minute the remaining time

for i in {25..1}
do
    echo "Remaining time: $i minutes"
    sleep 1m
done
