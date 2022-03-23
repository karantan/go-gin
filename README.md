# Go Gin
Playing around with go-gin web framework

In this example I'll show you how to create a simple API with JWT auth.

Why use tokens?
They are stateless, can be generated from anywhere, good access control.

## Usage

Install [Go](https://go.dev/doc/install) if you don't have it already.

```
$ git clone git@github.com:karantan/go-gin.git
$ cd go-gin
$ go run main.go
...
```

Open another terminal and make some API requests:

```
$ curl http://localhost:3000/
{"message":"Hello World"}

$ curl http://localhost:3000/api/v1/articles/
{"error":"No Authorization header provided"}

# get bearer JWT token
$ curl --header "Content-Type: application/json" --request POST --data '{"client_id": "admin", "client_secret": "secret"}' http://localhost:3000/login
{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFkbWluIiwiZXhwIjoxNjQ4MTI5ODY1fQ.GU8ZT57XMIrYYCbfUZUEg3c_brKvcwqo7oT8CtiTHRw","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsImV4cCI6MTY0ODY0ODI2NX0.3MMAmXHXoN0sIs52j_g59jP453G2hQQqTCQPi9gcDcA","token_type":"Bearer","expires_in":86400}

# make request again, this time with bearer token
$ curl --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFkbWluIiwiZXhwIjoxNjQ4MTI5ODY1fQ.GU8ZT57XMIrYYCbfUZUEg3c_brKvcwqo7oT8CtiTHRw" http://localhost:3000/api/v1/articles/
{"articles":[{"id":"1","user_id":100,"title":"Hi","slug":"hi"},{"id":"2","user_id":200,"title":"sup","slug":"sup"},{"id":"3","user_id":300,"title":"alo","slug":"alo"},{"id":"4","user_id":400,"title":"bonjour","slug":"bonjour"},{"id":"5","user_id":500,"title":"whats up","slug":"whats-up"}]}
```

## Limitations

Both article and user packages have mocked models (i.e. all data is hardcoded). There
is a boltdb database ready to be plugged in, but I haven't done it because I want to
keep this example minimalistic.

JWT authetication has also a few limitations:
- The JWT can only be invalidated when it expires
- The user will need to re-login after the token expires

We can solve these issues by storing JWT tokens in a database and implementing a token
refresh mechanism.

See this [post](https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/)
for an example.

## TODO

- [ ] Unit tests
- [ ] Add DB
