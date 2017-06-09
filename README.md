#WebRTC Messaging Server
##Author
Moad Ben-Suleiman

##Overview
A messaging server that complies with the [WebRTC standards] (https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling), including signal and ICE candidate exchanging

##How to run
Prerequisites: A running Go installation

If you have Go installed locally, simply run:
    go run *go
in the cloned repository and the server will start up on port :8080

To connect via websocket, the url should be in the following format:
    localhost:8080/ws?username={insertusernamehere}
usernames are passed in as a url query parameter.

###Testing
To test, I recommend installing wscat via npm and connecting using urls as listed above.

##Issues