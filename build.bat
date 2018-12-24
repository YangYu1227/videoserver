cd ./api
go build -o ../bin/api.exe
cd ..

cd ./streamserver
go build -o ../bin/streamserver.exe
cd ..

cd ./web
go build -o ../bin/web.exe
cd ..