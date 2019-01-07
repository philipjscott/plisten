# plisten

Register callback functions that trigger on tcp/udp requests. Currently only supports DNS; I'd be willing to add more functionality if people are interested :smiley:

## Usage

```
package main

import (
	"github.com/ScottyFillups/plisten/pkg/dnsl"
	"fmt"
)

func logDNSWarn(d *dnsl.DNSListener, match string) {
	fmt.Println("You visited: " + match + ". Shouldn't you be working?")
}

func main() {
	dataChan := make(chan dnsl.Packet)
	dl := dnsl.New()

	err := dl.Listen(dataChan)
	if err != nil {
		log.Fatal("Failed to initialize DNS listener")
	}

	err := dl.Register("*facebook*", logDNSWarn)
	if err != nil {
		log.Fatal("Failed to compile regexp")
	}

	for data := range dataChan {
		if data.Error != nil {
			fmt.Println(data.Error)
			dl.Close()
			break
		}

		fmt.Println(data.Host)
	}
}
```

See `/examples` for more specific and creative usages.
