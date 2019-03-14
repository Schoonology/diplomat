## Custom Script Test

```
> POST /post HTTP/1.1
> Content-Type: application/json
{{ custom }}
< HTTP/1.1 200 OK
{
  "args": {},
  "data": "this is custom\n",
  "files": {},
  "form": {},
  "headers": {
    "Content-Length": "15",
    "Content-Type": "application/json",
    "Host": "localhost:7357",
    "User-Agent": "Diplomat/0.0.1"
  },
  "json": null,
  "origin": "172.17.0.1",
  "url": "http://localhost:7357/post"
}
```
