# Authorization API

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=flat&logo=mongodb&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=flat&logo=swagger&logoColor=white)

> ⚠️ **Note:** THIS IS CURRENTLY A WORK IN PROGRESS, NOT ALL FEATURES HAVE BEEN IMPLEMENTED YET. THIS VERSION IS BUGGY AND WON'T NECESSARILY WORK.

## Overview
This is an **Authorization API** template that I use as a starting point for all my projects. It provides a production-ready authentication system built with **Golang** using the **Gin** framework. Feel free to use and customize it for your projects.

### Features
- **User and Session Management**
- **Built in Admin Dashboard with Analytics and Statistics**
- **MongoDB Integrated** and already setup
- **Salted Password Hashing**
- **Email Verification** already setup!
- **Multi-Factor Authentication (MFA) Support** via Email, but infrastructure allows for easy expansion.
- **Rate Limiting & Brute Force Protection** plus a lot of other security precautions
- **Swagger Documentation for Debugging** (Accessible at `/docs/index.html`)
- **and much much more!**

---

## Important Notes
Here are some important points you need to consider/know when using this template

### Liability Disclaimer

⚠️ This template is provided **as is**, without any warranties or guarantees. I am **not responsible or liable** for any issues, security vulnerabilities, data loss, or damages that may occur while using, modifying, or deploying this template. Use it at your own risk, and make sure to review and customize it according to your project's security and operational requirements.

### MongoDB setup
The API is currently set up to recieve requests from either a local MongoDB database instance or a remote one. Running it locally will require further installation and setup from the official MongoDB site. If you are using it remotely, for example Atlas, copy and paste your URI into the .env file (you may need to create one, the instruction for how to set it up is in the installation section. Docker will NOT install MongoDB, therefore you either need to download it yourself or create a remote instance.

### Debug mode
The API is set to debug mode by default and will need to be manually configured and removed before pushing to production.

---

## Setup

### Prerequisites

[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/products/docker-desktop/)
[![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)](https://www.mongodb.com/)

### Create an environment file
Create your .env file in the root directory of the auth-api folder.
Set it up as so:

```env
DATABASE_URI=<your database uri>
USER_AES_KEY=<your aes-256 encryption key>
RESEND_API_KEY=<your api key from resend>
ROOT_ADMIN_USER=<admin_username_here>
ROOT_ADMIN_PASSWORD=<admin_password_here>

BASE_URL=<the base url of your front end, e.g. example.com>
EMAIL_REDIRECT_BASE=<the base url for email redirects, e.g. example.com/email/redirect/token:>
```

Make sure you put your unique AES encryption key **INSIDE THE ENVIRONMENT FILE** and not inside the code. To generate a random secure key, visit [this website](https://generate-random.org/encryption-key-generator) - Leaving everything on default values.

For the ROOT ADMIN user section; this will be the account that you log in to that will manage all the other admin users. It will have the highest clearance level. Make sure NOT to share the login details of this account anywhere. Logging into the admin portal through this account is the only way you can create new admin users.

For the EMAIL REDIRECT BASE section; you will need to have a route defined on your front end which will allow a param to be placed after it continuing setup. For example, for ``test.com/v/:code`` - you would put ``test.com/v/`` as the env variable.

### Quick Start

```sh
# Build the docker images
make -f deploy/Makefile docker-build

# Run the docker images
make -f deploy/Makefile docker-run
```

This will run the worker process in the background, and start running the main API process in the current terminal.

To view the logs for the worker process, go to a different terminal and run:

```sh
docker logs -f auth-worker
```

---

## How it works  
This section will deep-dive into explaining every part of the API and how it works.

---

### Users
When a user is created the following will happen:

1. **Endpoint Called**
2. **User Password Hashed**
3. **Information sent to database for storage**

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

## Contributing

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

Feel free to contribute by submitting pull requests or reporting issues. Any suggestions for improvement are welcome!

---

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This project is open-source and available under the [MIT License](LICENSE).

---

<div align="center">
<sup>P.S. This is one of my first go projects so don't expect it to be too good.</sup>
</div>
