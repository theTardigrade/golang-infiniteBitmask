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

	fmt.Println("start:", start.String())
	fmt.Println("stop:", stop.String())
	fmt.Println("pause:", pause.String())
	fmt.Println("reset:", reset.String())

	startAndReset := g.ValueFromNames("start", "reset")

	fmt.Println("start & reset:", startAndReset.String())

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
