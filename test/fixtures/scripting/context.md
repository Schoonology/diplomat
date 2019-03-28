# Context

The `ctx` object provides a single, global "context" to store data between
requests.

## Setting a context

```
> GET /json
< 200 OK
{? chain(json.decode, ctx.set('body')) ?}
```

## Loading that context

```
> POST /post
> Content-Type: application/json
{{ json.encode(ctx.get('body')) }}
< 200 OK
{? chain(
  json.decode,
  get('json'),
  equal(ctx.get('body'))
) ?}
```
