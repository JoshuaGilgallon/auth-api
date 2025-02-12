# Authorisation API

## Overview
This is an **Authorization API** template that I use as a starting point for all my projects. It provides a production-ready authentication system built with **Golang** using the **Gin** framework. Feel free to use and customize it for your projects.

### Features
- **Database Integration**
- **Salted Password Hashing**
- **Multi-Factor Authentication (MFA)** via Email, SMS, and Authenticator Apps
- **Rate Limiting & Brute Force Protection**
- **Session Management**
- **Swagger Documentation** (Accessible at `/docs/index.html`)

---

## Installation

### Install Go
Download and install Go from the official site: [Go Installation Guide](https://go.dev/doc/install)

### Add Go's Bin Directory to PATH (macOS)
```sh
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
source ~/.bash_profile
```

### Install Swag (Swagger Documentation Generator)
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

#### Verify Swag Installation
```sh
swag -v
```
If you encounter issues, ensure Go is added to your system's PATH (see macOS steps above).

---

## Project Setup

For this part you can either run the Makefile to automatically do the following steps or do it yourself. Considering you will have to do this a large amount of times during development I recommend using and configuring the Makefile to your needs.

### Makefile
```sh
make -f deploy/Makefile run
```
If you are on windows ensure you have the Makefil file set to LF line endings, NOT CRLF.

### Initialize Swagger Documentation
```sh
swag init -g cmd/api/main.go
```

### Manage Dependencies
```sh
go mod tidy
```

### Run the API
```sh
go run cmd/api/main.go
```

---

## Contribution
Feel free to contribute by submitting pull requests or reporting issues. Any suggestions for improvement are welcome!

---

## License
This project is open-source and available under the [MIT License](LICENSE).

