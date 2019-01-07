package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
		t := time.Now().Format("15:04:05")

		if data.Error != nil {
			fmt.Fprintf(os.Stderr, "[%s] ERROR: %v\n", t, data.Error)
			dl.Close()
			break
		}

		fmt.Printf("[%s] HOST: %s\n", t, data.Host)
	}
}
