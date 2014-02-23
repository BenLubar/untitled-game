#!/bin/bash

set -e

GOOS=darwin   GOARCH=amd64  go build -o untitled-game-mac  "$@"
GOOS=linux    GOARCH=arm    go build -o untitled-game-arm  "$@"
GOOS=linux    GOARCH=amd64  go build -o untitled-game      "$@"
GOOS=windows  GOARCH=amd64  go build -o untitled-game.exe  "$@"

cd prototype/species
GOOS=darwin   GOARCH=amd64  go build -o species-mac  "$@"
GOOS=linux    GOARCH=arm    go build -o species-arm  "$@"
GOOS=linux    GOARCH=amd64  go build -o species      "$@"
GOOS=windows  GOARCH=amd64  go build -o species.exe  "$@"
cd ../..
