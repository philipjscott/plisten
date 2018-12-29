package main

import (
	"fmt"
	"log"

	"github.com/ScottyFillups/plisten/pkg/dnsl"
)

func logDns(d *dnsl.DNSListener, match string) {
	fmt.Println(match)
}

func handleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	dl := dnsl.New()

	err := dl.Register(".*", logDns)
	handleErr(err)

	err = dl.Listen()
	handleErr(err)
}
