#!/bin/bash

rm -f test.mov
ffmpeg -progress - -nostats -i pantaleo.mkv test.mov