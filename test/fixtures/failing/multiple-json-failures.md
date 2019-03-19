# Multiple JSON Failures

## Invalid JSON Schema

```
> GET /json HTTP/1.1
< HTTP/1.1 200 OK
<
{? json_schema(file("test/fixtures/failing/invalid-json-schema.json")) ?}
```

## Fail to Match JSON Schema

```
> GET /json HTTP/1.1
< HTTP/1.1 200 OK
<
{? json_schema(file("test/fixtures/failing/fail-to-match-schema.json")) ?}
```