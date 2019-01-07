package main

import (
	"fmt"
	"log"

	"github.com/ScottyFillups/plisten/pkg/dnsl"
)

func logDNS(d *dnsl.DNSListener, match string) {
	fmt.Println(match)
}

func handleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var dataChan chan dnsl.Packet
	dl := dnsl.New()

	err := dl.Register(".*", logDNS)
	handleErr(err)

	err = dl.Listen(dataChan)
	handleErr(err)

	fmt.Println("REE!!!")

	for data := range dataChan {
		fmt.Println("foo")

		if data.Error != nil {
			fmt.Println(data.Error)
			dl.Close()
			break
		}
	}
}
