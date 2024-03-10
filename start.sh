#!/bin/sh

# Start Envoy Proxy
envoy -c /etc/envoy.yaml &

# Start Server
./server