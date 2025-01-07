#!/bin/bash

port=${1:-$PORT}
id=${2:-"1"}

for i in {1..5}
do
    echo "[${id}] Test $i: $(curl -s http://0.0.0.0:${port}/ping)"
    sleep 10
done
