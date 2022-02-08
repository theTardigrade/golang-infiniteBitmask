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

	startBitmask := g.ValueFromName("start")
	stopBitmask := g.ValueFromName("stop")
	pauseBitmask := g.ValueFromName("pause")
	resetBitmask := g.ValueFromName("reset")

	fmt.Println("start:", startBitmask.Number().Text(2)) // 1
	fmt.Println("stop:", stopBitmask.Number().Text(2)) // 10
	fmt.Println("pause:", pauseBitmask.Number().Text(2)) // 100
	fmt.Println("reset:", resetBitmask.Number().Text(2)) // 1000

	startAndResetBitmask := g.ValueFromNames("start", "reset")

	if startAndResetBitmask.Number().Int64() == startBitmask.Number().Int64()|resetBitmask.Number().Int64() {
		fmt.Println("start and reset:", startAndResetBitmask.Number().Text(2)) // 1001
	}
}
```
