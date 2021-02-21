# create-google-jwt

A program that initiates an OAuth flow with Google and prints out the resultant JWT.

## Usage

```
go run main.go -clientID=xxx-yyy.apps.googleusercontent.com -clientSecret=zzzzzzzzzzzzzzzz
```

The program will then open up the Google authentication dialog in a browser. On completion, it will complete the handshake and print out the access and id tokens for use against APIs.
