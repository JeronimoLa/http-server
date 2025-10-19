# Chirp API (Go)

A simple RESTful API built in **Go** for managing chirps (messages) and tracking them against users. The API uses **JWT access tokens** for protected routes and **refresh tokens** to get new access tokens when they expire.


---

## Overview

- **Access protected routes** using jwt token in `Authorization` header.
- manage's the lifecycle of refresh tokens
- Add, update, and delete chirps 

## API Endpoints

### Auth Routes

| Method | Endpoint           | Description         | Auth |
| ------ | ------------------ | ------------------- | ---- |
| POST   | `/api/login`    | Log in and get access & refresh JWT | ❌    |
| POST   | `/api/refresh`    | Get a new access token using refresh token  | ✅    |
| POST   | `/api/revoke`    | Revoke refresh token  | ✅   |

### User Routes

| Method | Endpoint             | Description         | Auth |
| ------ | -------------------- | ------------------- | ---- |
| PUT    | `/api/users`      | Update user info       | ✅    |
| POST   | `/api/users`      | Create user    | ✅    |

### Admin Routes 

| Method | Endpoint             | Description         | Auth |
| ------ | -------------------- | ------------------- | ---- |
| GET    | `/admin/reset`      | Delete all user's along with data | ✅    |
| POST   | `/admin/metrics`      | Return number of visit's to app | ✅    |

### Chirp Routes
| Method | Endpoint             | Description         | Auth |
| ------ | -------------------- | ------------------- | ---- |
| GET    | `/api/chirps`      | Get all chirps | ✅    |
| POST    | `/api/chirps`      | Add chirp to user | ✅    |
| GET   | `/api/chirps/{chirpID}`      | Get single chirp that belongs to user | ✅    |
| DELETE | `/api/chirps/{chirpID}`      | Delete single chirp that belongs to user  | ✅    |
