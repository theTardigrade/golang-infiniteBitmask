# golang-infiniteBitmask

This Go package allows for the creation of bitmasks of theoretically infinite size.

## Example

```golang
package main

import (
	"fmt"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

func main() {
	g := bitmask.NewGenerator()

	start := g.ValueFromName("start")
	stop := g.ValueFromName("stop")
	pause := g.ValueFromName("pause")
	reset := g.ValueFromName("reset")

	fmt.Println("start:", start.String()) // 1
	fmt.Println("stop:", stop.String())   // 10
	fmt.Println("pause:", pause.String()) // 100
	fmt.Println("reset:", reset.String()) // 1000

	startAndReset := g.ValueFromNames("start", "reset")

	fmt.Println("start and reset:", startAndReset.String()) // 1001

	if startAndReset.Number().Int64() == start.Number().Int64()|reset.Number().Int64() {
		fmt.Println("match [1]")
	}

	startAndReset2 := start.Clone()
	startAndReset2.Combine(reset)

	if startAndReset.Equal(startAndReset2) {
		fmt.Println("match [2]")
	}

	if startAndReset.String() == startAndReset2.String() {
		fmt.Println("match [3]")
	}
}
```
