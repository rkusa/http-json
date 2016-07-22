# http-json

JSON helpers to read (both filtered and unfiltered) and write JSON request bodies.

[![Build Status][travis]](https://travis-ci.org/rkusa/http-json)
[![GoDoc][godoc]](https://godoc.org/github.com/rkusa/http-json)

### Example

```go
json.Read(r, &user)
json.ReadFiltered(r, &user, []string{"name"})
json.Write(rw, &user)
```

## License

[MIT](LICENSE)

[travis]: https://img.shields.io/travis/rkusa/http-json.svg
[godoc]: http://img.shields.io/badge/godoc-reference-blue.svg
