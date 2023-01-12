# Vue WebSocket Test

## Introduction

Learning about WebSockets with Vue and Go.

It is a small chat app, with only one channel so far, which shows messages
and displays active users in real-time. 

So far it uses websockets for all data. WIP is making it use an API,
which will fire events through the WebSocket. It'll also be adding 
user login/registration and JWTs. Maybe different channels or group
chats will be added later.

## Running

Requires:
- Golang installed
- Nodejs and npm
- You'll likely need a monitor to try it out too

### Backend

Runs on port :8080 by default.

```
cd api/cmd/api
go build
./api
```

### Frontend

```
cd web
npm install
npm run dev
```
