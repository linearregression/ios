#!/bin/bash

# testing latency for 3 to 9 servers with a single client
for (( i = 3; i <= 9; i = i+2 )); do
	./scripts/simple_test.sh $i 1
done