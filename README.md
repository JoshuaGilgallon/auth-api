# THIS IS CURRENTLY A WORK IN PROGRESS, NOT ALL FEATURES HAVE BEEN IMPLEMENTED YET. THIS VERSION IS BUGGY AND WON'T NECESSARILY WORK. 

---

# Authorisation API

## Overview
This is an **Authorization API** template that I use as a starting point for all my projects. It provides a production-ready authentication system built with **Golang** using the **Gin** framework. Feel free to use and customize it for your projects.

### Features
- **User Management**
- **Session Management**
- **MongoDB Integrated** and already setup
- **Salted Password Hashing**
- **User Data AES-256 Encryption**
- **Multi-Factor Authentication (MFA)** via Email, SMS, and Authenticator Apps
- **Rate Limiting & Brute Force Protection**
- **Swagger Documentation** (Accessible at `/docs/index.html`)
- **and much much more!**


---

## Important Notes
Here are some important points you need to consider/know when using this template

### Liability Disclaimer

This template is provided **as is**, without any warranties or guarantees. I am **not responsible or liable** for any issues, security vulnerabilities, data loss, or damages that may occur while using, modifying, or deploying this template. Use it at your own risk, and make sure to review and customize it according to your project’s security and operational requirements.

### MongoDB setup
The API is currently set up to recieve requests from either a local MongoDB database instance or a remote one. Running it locally will require further installation and setup from the official MongoDB site. If you are using it remotely, for example Atlas, copy and paste your URI into the .env file (you may need to create one, the instruction for how to set it up is in the installation section.

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

### Create a .env file
Create your .env file in the root directory of the auth-api folder.
Set it up as so:

```env
DATABASE_URI=<your database uri>
USER_AES_KEY=<your aes-256 encryption key>
ROOT_ADMIN_USER=<admin_username_here>
ROOT_ADMIN_PASSWORD=<admin_password_here>
```

Make sure you put your unique AES encryption key **INSIDE THE ENVIRONMENT FILE** and not inside the code. To generate a random secure key, visit [this website](https://generate-random.org/encryption-key-generator) - Leaving everything on default values.

For the ROOT ADMIN user, this will be the account that you log in to that will manage all the other admin users. It will have the highest clearance level. Make sure NOT to share the login details of this account anywhere. Logging into the admin portal through this account is the only way you can create new admin users.

## Project Setup
For this part you can either run the Makefile to automatically do the following steps or do it yourself. Considering you will have to do this a large amount of times during development I recommend using and configuring the Makefile to your needs.

### Makefile
```sh
make -f deploy/Makefile run
```
If you are on windows ensure you have the Makefile file set to LF line endings, NOT CRLF.

## Manual Setup

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

OR build it
```sh
go build -o bin/api ./cmd/api/main.go
./bin/api
```

---

## How it works  
This section will deep-dive into explaining every part of the API and how it works.

---

### Users
When a user is created the following will happen:


1. **Endpoint Called**
2. **User Password Hashed**
3. **User Email and Phone number encrypted with AES-256**
4. **Information sent to database for storage**

When a user is fetched the following will happen

1. **Endpoint Called**
2. **User Email and Phone Number decrypted with AES-256**
3. **Information displayed** 

---

### Session Handling

Session handling ensures users stay logged in while keeping their accounts secure. Our implementation follows **best practices using Access Tokens and Refresh Tokens** for authentication.  

#### **How It Works**  

1. **User Logs In**  
   - The server verifies the user's credentials.  
   - It returns **two tokens**:  
     - **Access Token** – Short-lived (30 min), used for authentication in API requests.  
     - **Refresh Token** – Long-lived (7 days), stored in an **HTTP-only Secure Cookie**.  

2. **Access Token Usage**  
   - The client stores the **Access Token in memory** (e.g., React state, Vuex, Redux).  
   - All API requests include:  
     ```
     Authorization: Bearer <access_token>
     ```
   - The server **validates** the token before processing the request.  

3. **Token Expiry & Renewal**  
   - When the **Access Token expires (after 30 min)**, API requests fail with `401 Unauthorized`.  
   - The client **automatically requests a new Access Token** using the **Refresh Token**.  
   - The browser **sends the Refresh Token in an HTTP-only cookie** to `/refresh`.  
   - If the Refresh Token is valid, the server issues a **new Access Token**.  

4. **Session Expiry & Logout**  
   - If the Refresh Token is expired (after 7 days) or revoked, the user is logged out.  
   - The client must **log in again** to get a new session.

---

## Contribution
Feel free to contribute by submitting pull requests or reporting issues. Any suggestions for improvement are welcome!

---

## License
This project is open-source and available under the [MIT License](LICENSE).

---

P.S. This is one of my first go projects so don't expect it to be too good.
