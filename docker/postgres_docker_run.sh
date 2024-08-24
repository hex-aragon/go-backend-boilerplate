#!/bin/bash

echo "docker run"
docker run --name go-boiler-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=boiler -d postgres:latest

