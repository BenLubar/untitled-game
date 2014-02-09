#!/bin/bash

set -e

GOOS=linux    GOARCH=arm    go build -o untitled-game-arm
GOOS=linux    GOARCH=amd64  go build -o untitled-game
GOOS=windows  GOARCH=amd64  go build -o untitled-game.exe
