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
	visits++

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
	var dataChan chan dnsl.Packet
	dl := dnsl.New()

	err := dl.Register(".*googlevideo.*", shutdown)
	handleErr(err)

	err = dl.Listen(dataChan)
	handleErr(err)

	for data := range dataChan {
		if data.Error != nil {
			fmt.Println(data.Error)
			dl.Close()
			break
		}
	}
}
