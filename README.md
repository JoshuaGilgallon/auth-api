# Authorisation API
This is an authrosation API that I use as a starter template for all my projects.
Feel free to use it too. It is written in golang with the gin API framework.

This repository has database integration, hashing, and everything you would need for a production level auth API.

Swagger is also already setup using swaggo/swag. Access it at /docs/index.html


## Installing go
Download go from [here](https://go.dev/doc/install)

## ON MACOS - Add go's bin directory to path
```echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile```
<br>
then
<br>
``source ~/.bash_profile``

## Installing SWAG cmd:
``go install github.com/swaggo/swag/cmd/swag@latest``

verify swag installation with:

``swag -v``

## Init SWAG
``swag init -g cmd/api/main.go``

## Go tidy
``go tidy``

## Run API
``go run cmd/api/main.go``