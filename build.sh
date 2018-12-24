echo "start!"

cd ~/desktop/gowork/src/video_server/api
go build -o ~/desktop/gowork/src/video_server/bin/api

cd ~/desktop/gowork/src/video_server/streamserver
go build -o ~/desktop/gowork/src/video_server/bin/streamserver

cd ~/desktop/gowork/src/video_server/web
go build -o ~/desktop/gowork/src/video_server/bin/web

echo "finished!"