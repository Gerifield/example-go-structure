# Example Go project structure

# First example

```
$ go run src/cmd/a/a.go
```


# Second example

```
$ go run src/cmd/b/b.go
```

Test the main endpoints with curl:

```
$ curl http://127.0.0.1:8080/
Hello b!
```

Test the secret endpoint without credentials (I like to use verbose mode to see the headers):

```
$ curl -v http://127.0.0.1:8080/secret
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET /secret HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Mon, 21 Dec 2020 16:42:53 GMT
< Content-Length: 0
<
* Connection #0 to host 127.0.0.1 left intact
* Closing connection 0
```

Now use the correct creds:

```
$ curl -v -u "admin:admin" http://127.0.0.1:8080/secret
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
* Server auth using Basic with user 'admin'
> GET /secret HTTP/1.1
> Host: 127.0.0.1:8080
> Authorization: Basic YWRtaW46YWRtaW4=
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Mon, 21 Dec 2020 16:53:48 GMT
< Content-Length: 17
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 127.0.0.1 left intact
This is a secret!* Closing connection 0
```

The second middleware should also response to and invalid password:

```
$ curl -v -u "test1:test1a" http://127.0.0.1:8080/secret
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
* Server auth using Basic with user 'test1'
> GET /secret HTTP/1.1
> Host: 127.0.0.1:8080
> Authorization: Basic dGVzdDE6dGVzdDFh
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 401 Unauthorized
< Date: Mon, 21 Dec 2020 16:54:10 GMT
< Content-Length: 5
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 127.0.0.1 left intact
Nope!* Closing connection 0
```

Third example with some JSON parsing:

```
$ curl http://127.0.0.1:8080/json -d '{"A":"test"}'
{"A":"test","B":"some value"}
```


# Third example

Surprisingly it is in the `example2` folder. It shows an example how to separate the business logic from the transport layer.

Run it:

```
$ go run src/cmd/c/c.go
```

Test calls:

```
$ curl 127.0.0.1:8080
hello index
```

```
$ curl "127.0.0.1:8080/greeting?name=Bob"
hello Bob
```
