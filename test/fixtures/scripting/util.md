# Provided utilities

## Teeing data

The `tee_file` method can be used to save response data to disk for use in other treaties. This is useful, for example, for access tokens. You only need log in once, in a treaty you run before the other treaties.

```
> GET /json
< 200 OK
{? tee_file('json.json') ?}
```

That data can be retrieved during request creation with the `file` function.

```
> POST /post
> Content-Type: application/json
{{ file('json.json') }}
< 200 OK
{? chain(
  json.decode,
  get('json'),
  get('slideshow'),
  get('author'),
  equal('Yours Truly')
) ?}
```
