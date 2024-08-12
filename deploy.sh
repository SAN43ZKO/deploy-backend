#!/bin/bash
docke compose stop
sleep 10
docker compose up --build -d