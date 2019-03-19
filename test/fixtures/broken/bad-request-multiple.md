# Multiple Bad Requests

## Bad Request Header

```
> POST /post HTTP/1.1
> Accept: text/plain
> Content-Type
The request body
< HTTP/1.1 200 OK
{}
```

## Bad Request Line

```
> INVALID
> Accept: text/plain
> Content-Type: application/json
The request body
< HTTP/1.1 200 OK
{}
```