#!/bin/bash

# Measure CPU utilization every 0.5 seconds
while true
do
    cpu_usage=$(top -bn2 -d 0.5 | fgrep 'Cpu(s)' | tail -1 | awk  -F'id,' '{ n=split($1, vals, ","); v=vals[n]; sub("%", "", v); printf "%f", 100 - v }')
    echo $cpu_usage
done
