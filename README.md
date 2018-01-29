# Overview

This is a Go implementation of the [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm used as a rate limiter.

It does not use what might be considered a more Go-like approach of a goroutine with a time.Ticker sending over a channel, but this keeps the implementation simpler.

This implementation is thread safe, and the one operation defined on the struct is protected by a sync.Mutex.

# Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/adamdrake/tokenbucket"
)

func main() {
    b := tokenbucket.NewTokenBucket(uint64(10), uint64(10), 3 * time.Second)
    i := 0
    for _ = range time.NewTicker(1 * time.Second).C {
        fmt.Println(i, b.Take())
        i++
    }
}
```