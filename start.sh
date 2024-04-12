#!/bin/bash

# turn on bash's job control
set -m

# Start Envoy Proxy
envoy -c /etc/envoy.yaml &

# Start Server
./server

fg %1