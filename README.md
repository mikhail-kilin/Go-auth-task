# Go-auth-task
##Deployed to
https://golang-auth-task.herokuapp.com/

## Setup
1) Clone this repository
2) Open project and use command: bin/setup
3) Set environment variables in .env

## Start
go run main.go

## Requests
### Sugn Up
POST /signup
{
  "email" : "email@email.com",
  "name" : "Name",
  "password" : "password"
}

### Sign in

POST /login
{
  "email" : "email@email.com",
  "password" : "password"
}

### Profile
GET /profile
Headers: Access-token, Refresh-token

### Refresh
POST /refresh
Headers: Access-token, Refresh-token

### Delete Refresh token
POST /refresh/delete
Headers: Access-token, Refresh-token

### Delete All Refresh tokens
POST /refresh/all/delete
Headers: Access-token, Refresh-token
