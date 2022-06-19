#!/bin/bash

echo "0" > /dev/stderr
for i in 20 40 60 80 100; do
    sleep 2
    echo "$i" > /dev/stderr
done;

date