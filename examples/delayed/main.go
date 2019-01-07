package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/ScottyFillups/plisten/pkg/dnsl"
)

const (
	limit = 2
	delay = 5
)

var (
	visits = 0
)

func scheduleShutdown() {
	time.Sleep(time.Minute * delay)
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func shutdown(d *dnsl.DNSListener, match string) {
	visits++

	fmt.Printf("%d requests remaining until delayed shutdown\n", limit-visits)

	if visits >= limit {
		fmt.Printf("Shutting down in %d minutes\n", delay)
		go scheduleShutdown()
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
			fmt.Println(data.Error )
			dl.Close()
			break
		}
	}
}
