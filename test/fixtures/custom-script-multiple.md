## Multiple Custom Script Test

```
> GET /response-headers?Custom={{header}} HTTP/1.1
< HTTP/1.1 200 OK
< Custom: {{ header }}
```
