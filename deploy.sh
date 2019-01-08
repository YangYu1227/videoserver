#! /bin/bash

mkdir ./bin/upload

cd bin

nohup ./api &
nohup ./web &

echo "deploy finished"
