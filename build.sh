#! /bin/bash

if [[ $1 == "pi" ]]
then
    echo "Building for Rasbian";
    env GOOS=linux GOARCH=arm GOARM=5 go build -o AutoDrone ./AD_Main;
    echo "Done.";
fi
if [[ $1 == "w10" ]]
then
    echo "Building for Windows 10 OS";
    go build -o AutoDrone.exe ./AD_Main;
    echo "Done.";
fi
if [[  $1 == "" ]]
then
    echo "Usage: { pi (Raspberry Pi) | w10 (Windows 10) }";
fi