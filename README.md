# WebRTC Messaging Server
## Author
Moad Ben-Suleiman

## Overview
A messaging server that complies with the [WebRTC standards](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling), including signal and ICE candidate exchanging

## How to run
Prerequisites: A running Go installation

If you have Go installed locally, simply run:
```
go run *.go
```
in the cloned repository and the server will start up on port :8080

To connect via websocket, the url should be in the following format:
```
localhost:8080/ws?username={insertusernamehere}
```
Usernames are passed in as a url query parameter.

### Testing
To test, I recommend installing wscat via npm and connecting using urls as listed above.

## Issues
1. Currently, the setting of remote descriptions occurs after an offer successfully reaches its destinations, this should instead be handled in the front-end client when a user accepts a connection
2. Promises aren't currently in use.