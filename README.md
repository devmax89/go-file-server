# go-file-server

Build executable:
go build -ldflags "-X main.Version=1.0.0" -o web-server

Build image:
docker build -t web-server:1.0.0 .