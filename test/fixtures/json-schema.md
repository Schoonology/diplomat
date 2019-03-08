# JSON Schema Test

```
> GET /json HTTP/1.1
< HTTP/1.1 200 OK
<
{? json_schema(file("test/fixtures/schema.json")) ?}
```
