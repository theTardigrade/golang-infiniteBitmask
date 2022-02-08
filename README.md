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

	if startAndReset.Number().Uint64() == start.Number().Uint64()|reset.Number().Uint64() {
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

	empty := start.Clone()
	empty.Clear()

	if empty.Number().Uint64() == 0 {
		fmt.Println("match [4]")
	}

	fmt.Println(g.String()) // ["start","stop","pause","reset"]
}
```
