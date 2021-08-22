#!/usr/bin/env bash

# Create data folders.
mkdir -p data/isaac
mkdir -p data/isaac/log
mkdir -p data/mysql
mkdir -p data/mysql/db
mkdir -p data/mysql/config

# Launch service. 
docker-compose up -d