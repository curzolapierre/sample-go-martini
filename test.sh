#!/bin/bash

# Set the port, could be passed as first argument or default to env variable PORT
port=${1:-$PORT}

for i in {1..18}
do
    echo "Test $i: $(curl -s http://0.0.0.0:${port}/ping)"
    sleep 10
done
