package main

import (
	"github.com/ScottyFillups/dnscap/pkg/dnscap"
	"fmt"
)

func logDns(d *dnscap.DNSCapturer, match string) {
	fmt.Println(match)
}

func logDnsPanic(d *dnscap.DNSCapturer, match string) {
	fmt.Println("BOT NET ALERT!!!", match)
}


func main() {
	cap := dnscap.New()

	cap.Register(".*", logDns)
	cap.Register(".*microsoft.*", logDnsPanic)

	cap.Listen()
}