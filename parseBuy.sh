#!/bin/sh
for n in `seq 0 $1`
do
    y=$((129 + n * 21))
    xdotool search --name 'EverQuest' mousemove 1100 $y click 1 sleep .5
done
