#!/bin/sh

apk add git
apk add curl

curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

GOPRIVATE=gitlab.com

go mod tidy

air