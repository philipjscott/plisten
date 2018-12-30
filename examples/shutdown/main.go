package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/ScottyFillups/plisten/pkg/dnsl"
)

const limit = 5

var (
	visits = 0
)

func shutdown(d *dnsl.DNSListener, match string) {
	visits += 1

	fmt.Printf("%d requests remaining until shutdown\n", limit-visits+1)

	if visits > limit {
		syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	}
}

func handleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var errChan chan error
	dl := dnsl.New()

	err := dl.Register(".*googlevideo.*", shutdown)
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
