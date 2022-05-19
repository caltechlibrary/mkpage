package mkpage

import (
	"fmt"
	"os"
	"sort"
)

// ShowEnvironent writes the contents of the environment to stdout.
// It is implemented so I can see the environment when mkpage is
// installed as a snap or in a container.
func ShowEnvironment() {
	keys := os.Environ()
	sort.Strings(keys)
	for i, key := range keys {
		fmt.Printf("%d %s\n", i, key)
	}
}
