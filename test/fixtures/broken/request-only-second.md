## First: Correct

```
> GET /status/200 HTTP/1.1
>
< HTTP/1.1 200 OK
< Content-Length: 0
<
```

## Second: Wrong

```
> GET /status/200 HTTP/1.1
>
```
