#!/bin/bash -xe
docker-compose down
docker system prune –a
docker system prune -f
