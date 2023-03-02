[![Go Reference](https://pkg.go.dev/badge/github.com/GreenLightning/go-lineiter.svg)](https://pkg.go.dev/github.com/GreenLightning/go-lineiter)

This package provides an allocation-free, zero-copy line iterator.

```go
    it := MakeLineIterator(...)
    for it.Next() {
        var line []byte = it.Bytes()
        fmt.Printf("%s\n", line)

        // More convenient, but allocates a string:
        // fmt.Printf("%s\n", it.Text())
    }
```

Semantically, this is equivalent to the following code, except that a
carriage return before the newline is trimmed by the line iterator:

```go
    for _, line := range strings.Split(..., "\n") {
        fmt.Printf("%s\n", line)
    }
```
