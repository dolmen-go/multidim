# multidim: Go multidimensional slices

[![Go Reference](https://pkg.go.dev/badge/github.com/dolmen-go/multidim.svg)](https://pkg.go.dev/github.com/dolmen-go/multidim)
[![Travis-CI](https://api.travis-ci.org/dolmen-go/multidim.svg?branch=main)](https://travis-ci.org/dolmen-go/multidim)
[![Go Report Card](https://goreportcard.com/badge/github.com/dolmen-go/multidim)](https://goreportcard.com/report/github.com/dolmen-go/multidim)


```go
import "github.com/dolmen-go/multidim"
```

## Examples

See also [examples in the documentation](https://pkg.go.dev/github.com/dolmen-go/multidim#pkg-examples).

Allocate a 2x2 matrix of strings with value `"foo"` in each cell ([run on the Go Playground](https://play.golang.org/p/V0pKeWGAy7z)):

```go
var square [][]string
multidim.Init(&square, "foo", 2, 2)
```