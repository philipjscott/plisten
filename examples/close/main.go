package main

import (
	"fmt"
	"log"

	"github.com/ScottyFillups/plisten/pkg/dnsl"
)

func closeSniffer(d *dnsl.DNSListener, match string) {
	fmt.Println("You visited: " + match + ". Stopping sniffer!")

	d.Close()
}

func handleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var errChan chan error
	dl := dnsl.New()

	err := dl.Register(".*test.*", closeSniffer)
	handleErr(err)

	err = dl.Listen(errChan)
	handleErr(err)

	for err := range errChan {
		if err != nil {
			fmt.Println(err)
			dl.Close()
			break
		}
	}
}
