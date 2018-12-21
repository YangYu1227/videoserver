cd D:/work/videoserver/api
env GOOS=linux GOARCH=amd64 go build -o D:/work/videoserver/bin/api

cd D:/work/videoserver/scheduler
env GOOS=linux GOARCH=amd64 go build -o D:/work/videoserver/bin/scheduler

cd D:/work/videoserver/streamserver
env GOOS=linux GOARCH=amd64 go build -o D:/work/videoserver/bin/streamserver

cd D:/work/videoserver/web
env GOOS=linux GOARCH=amd64 go build -o D:/work/videoserver/bin/web