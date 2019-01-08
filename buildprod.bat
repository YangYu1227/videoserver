cd ./api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api.exe
cd ..

cd ./web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web.exe
cd ..