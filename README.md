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
	cap := dnscap.New()

	cap.Register(".*", logDns)
	cap.Register(".*facebook.*", logDnsWarn)

	cap.Listen()
}
```

See `/examples` for more creative usages.
