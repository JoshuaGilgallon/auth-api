# Authorisation API

## Installing SWAG cmd:
``go install github.com/swaggo/swag/cmd/swag@latest``

## ON MACOS - Add go's bin directory to path
```echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile```
<br>
then
<br>
``source ~/.bash_profile``

verify swag installation with:

``swag -v``

## Init SWAG
``swag init -g cmd/api/main.go``