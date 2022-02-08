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

	fmt.Println("start:", start.Number().Text(2))
	fmt.Println("stop:", stop.Number().Text(2))
	fmt.Println("pause:", pause.Number().Text(2))
	fmt.Println("reset:", reset.Number().Text(2))

	startAndReset := g.ValueFromNames("start", "reset")

	if startAndReset.Number().Int64() == start.Number().Int64()|reset.Number().Int64() {
		fmt.Println("start & reset:", startAndReset.Number().Text(2))
	}
}

```
