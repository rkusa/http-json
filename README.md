# http-json

JSON helper to read JSON request bodies filtering the read JSON properties by a provided whitelist.

[![GoDoc][godoc]](https://godoc.org/github.com/rkusa/http-json)

### Example

```go
json.Read(req, &user, []string{"name"})
```

[godoc]: http://img.shields.io/badge/godoc-reference-blue.svg