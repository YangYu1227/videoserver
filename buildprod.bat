cd ./api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api
cd ..

cd ./streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver
cd ..

cd ./web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web