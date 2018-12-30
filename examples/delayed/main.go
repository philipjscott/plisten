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
	visits += 1

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
