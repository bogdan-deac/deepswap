# deepswap

This package allows you to swap values in deeply nested data structures.

```go
package main

import (
  "fmt"

  "github.com/bogdan-deac/deepswap"
)

type Ugly = struct {
  A string
  B int
  C *int
  D []int
  E map[string]int
}

func main() {
  uggo := Ugly{
    A: "1",
    B: 1,
    C: &[]int{1}[0],
    D: []int{1, 1},
    E: map[string]int{
      "1": 1,
    },
  }
  fmt.Println(deepswap.DeepSwap(uggo, 1, 3))
}
```
