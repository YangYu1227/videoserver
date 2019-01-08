#! /bin/bash

cd bin

nohup ./api &
nohup ./web &

echo "deploy finished"
