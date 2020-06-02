# go-chat
====

Chat server and client with golang based on socket.io lib

- - - - 

## Necessary Technology Versions

Technology  | Version
------------- | -------------
Go | go1.14.3 linux/amd64
Docker | 18.09.6
docker-compose | 1.24.1

## Running Chat Server

To run the chat server we create a docker container for it

    $ docker-compose up -d

## Running Chat Client

To run the chat client we will run it locally

    $ go run cmd/client/client.go

## Configurations

### Server Environment Variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| SERVER_NAME | Server Name | "" |
| PORT | Server Port | 3000 |

### Client Environment Variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| SERVER_IP | Server IP | localhost |
| SERVER_PORT | Server Port | 3000 |