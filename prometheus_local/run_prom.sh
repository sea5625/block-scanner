#!/bin/sh

# Create network.
docker network rm prom_net
docker network create prom_net

# Remove used containers.
docker rm -f prometheus
docker rm -f loopchain_export

# Launch prometheus.
docker run -d -p 9090:9090 -v ${PWD}/conf:/prometheus-data \
--network prom_net \
--name=prometheus \
prom/prometheus --config.file=/prometheus-data/prometheus.yml \

# Launch loopchain_exporter from docker image.

docker run -d -p 9095:9095 -v ${PWD}/conf:/conf -e INTERVAL=5 -e TIMEOUT=2  \
--network prom_net \
--link prometheus \
--name=loopchain_export \
iconloop/loopchain_exporter:0.0.8a
