#!/bin/bash -xe
echo "||||||||||||||||||||||||||||||||||||||||||||"
echo "||        Starting Local Development      ||"
echo "||               Arthur-Web               ||"
echo "||||||||||||||||||||||||||||||||||||||||||||"

echo "Launching Infrastructure"
echo ""
echo "==> Launch Docker"
open --background -a Docker
while ! docker system info > /dev/null 2>&1; do sleep 1; done
echo "docker running"
echo ""
docker network create arthur
echo "==> Launch app via docker compose"
docker-compose up -d --build && docker-compose logs -f
echo ""

