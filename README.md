# Authorisation API

# Installing go
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
``go run -g cmd/api/main.go``