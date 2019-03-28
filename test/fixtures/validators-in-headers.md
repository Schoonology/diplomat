## Date headers

Date headers are, by their nature, both required and never identical. Instead,
we can use the custom `date` validator.

```
> GET /get
< 200 OK
< Date: {? is_date ?}
```

## Regex-testing cookie contents

Since Cookie headers are comprised of multiple pieces, it's a good idea to use
a regex to test their contents.

```
> GET /cookies/set/key/value
< 302 FOUND
< Set-Cookie: {? regexp("key=value") ?}
```
