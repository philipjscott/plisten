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
	var dataChan chan dnsl.Packet
	dl := dnsl.New()

	err := dl.Register(".*test.*", closeSniffer)
	handleErr(err)

	err = dl.Listen(dataChan)
	handleErr(err)

	for data := range dataChan {
		if data.Error != nil {
			fmt.Println(data.Error )
			dl.Close()
			break
		}
	}
}
