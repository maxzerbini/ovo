#!/bin/bash

./ovo -conf=./conf/serverconf.json > node_1.log 2>&1 &
sleep 3
./ovo -conf=./conf/serverconf2.json > node_2.log 2>&1 &
sleep 3
./ovo -conf=./conf/serverconf3.json > node_3.log 2>&1 &