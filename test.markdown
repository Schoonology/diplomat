# Markdown Test

## Markdown: GET /status/422 

```
> GET /status/422 HTTP/1.1
>
< HTTP/1.1 422 UNPROCESSABLE ENTITY
< Content-Length: 0
```

## Markdown: Get /status/200

```
> GET /status/200 HTTP/1.1
>
< HTTP/1.1 200 OK
< Content-Length: 0
```

## Markdown: Get Response

```
> GET /base64/RGlwbG9tYXQgaXMgYXdlc29tZSEK HTTP/1.1
>
< HTTP/1.1 200 OK
<
Diplomat is awesome!
```