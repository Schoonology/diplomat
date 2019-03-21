# Output

## Test: Passing

```
> GET /status/200 HTTP/1.1
>
< HTTP/1.1 200 OK
< Content-Length: 0
<
```

## Test: Failing Assertion

```
> GET /status/200 HTTP/1.1
>
< HTTP/1.1 422 UNPROCESSABLE ENTITY
< Content-Length: 0
<
```

## Test: Failing due to Bad Template

```
> POST /post HTTP/1.1
> Accept: text/plain
> Content-Type
The request body
< HTTP/1.1 200 OK
{}
```