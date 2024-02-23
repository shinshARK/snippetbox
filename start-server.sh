#! /bin/bash
# go run ./cmd/web >>./tmp/info.temp.log 2>>./tmp/error.temp.log
go run ./cmd/web 2> >(tee -a ./tmp/error.log >&2) | tee -a ./tmp/info.log
