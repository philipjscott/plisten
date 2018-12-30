# plisten

Register callback functions that trigger on tcp/udp requests. Currently only supports DNS; I'd be willing to add more functionality if people are interested :smiley:

## Usage

```
package main

import (
	"github.com/ScottyFillups/plisten/pkg/dnsl"
	"fmt"
)

func logDns(d *dnsl.DNSListener, match string) {
	fmt.Println(match)
}

func logDnsWarn(d *dnsl.DNSListener, match string) {
	fmt.Println("You visited: " + match + ". Shouldn't you be working?")
}

func main() {
	dl := dnsl.New()

	dl.Register(".*", logDns)
	dl.Register(".*facebook.*", logDnsWarn)
	dl.Listen()
}
```

See `/examples` for more specific and creative usages.
