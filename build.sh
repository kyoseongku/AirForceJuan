#! /bin/bash

echo "Building for Raspberry Pi";

env GOOS=linux GOARCH=arm GOARM=5 go build -o AutoDrone.exe ./AD_Main;

echo "Done";
