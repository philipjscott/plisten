package main

import (
	"fmt"
	"log"

	"github.com/ScottyFillups/plisten/pkg/dnsl"
)

func handleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	dataChan := make(chan dnsl.Packet)
	dl := dnsl.New()

	err := dl.Listen(dataChan)
	handleErr(err)

	for data := range dataChan {
		if data.Error != nil {
			fmt.Println(data.Error)
			dl.Close()
			break
		}

		fmt.Println(data.Host)
	}
}
