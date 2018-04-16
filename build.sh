#! /bin/bash

echo "Building for Raspberry Pi";

env GOOS=linux GOARCH=arm GOARM=5 go build;

echo "Done";
