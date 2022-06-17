#!/bin/bash

docker compose up -d --build
JWT_SECRET_KEY=dummy-secret-key JWT_SECONDS_TO_EXPIRE=86400 go test -v ./... && docker compose down
