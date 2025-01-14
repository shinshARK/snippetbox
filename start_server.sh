#!/bin/bash

# Default values for flags
ADDR=":4000"  # Default server address

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -addr=*)
      ADDR="${1#*=}"
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Start the Go server with the parsed address
go run ./cmd/web -addr="$ADDR" 2>&1 | tee -a ./tmp/logs.log