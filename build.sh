#! /bin/bash

if [[ $1 == "pi" ]]
then
    echo "Building for Rasbian";
    env GOOS=linux GOARCH=arm GOARM=5 go build -o AutoDrone;
    echo "Done.";
elif [[ $1 == "w10" ]]
then
    echo "Building for Windows 10 OS";
    go build -o AutoDrone.exe;
    echo "Done.";
else
    echo "Usage: ./build.sh [ pi (Raspberry Pi) | w10 (Windows 10) ]";
fi
