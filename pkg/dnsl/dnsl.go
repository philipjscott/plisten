// Package dnsl listens for packets and decodes DNS information;
// it provides a simple API for registering callbacks that are called on a Regexp match:
//
//  func main() {
//  	dl := dnsl.New()
//  	dl.Register(".*foobar.*")
//
//  	// You can pass in a channel to receive error information
//  	dl.Listen(nil)
//
//  	// Infinite loop to prevent program from exiting; in "real" usage you'll
//  	// likely be iterating over the error channel, and calling dl.Close() on error,
//  	// or some other condition
//  	for {}
//  }
package dnsl

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ScottyFillups/plisten/internal/dnscap"
)

// DNSListener contains methods for listening for DNS requests
type DNSListener struct {
	handlers map[*regexp.Regexp]func(*DNSListener, string)
	active   bool
}

// Packet stores information about the DNS request packet
type Packet struct {
	// Host is the requested hostname
	Host string
	// Error is non-nil if an error occurred when decoding the packet
	Error error
}

// New returns a DNSListener struct.
func New() DNSListener {
	return DNSListener{
		handlers: make(map[*regexp.Regexp]func(*DNSListener, string)),
		active:   false,
	}
}

// Register maps a regular expression to a callback function.
// The callback function is triggered if the DNS question matches the regular expression.
func (d *DNSListener) Register(regexStr string, handler func(*DNSListener, string)) error {
	regex := regexp.MustCompile(regexStr)

	if regex == nil {
		return errors.New("Failed to compile regular expression")
	}

	d.handlers[regex] = handler

	return nil
}

// Close stops listening for packets
func (d *DNSListener) Close() {
	d.active = false
}

// Listen begins listening for packets on all active internet connections (promiscuous mode is disabled).
// Returns an error if the packet capturer fails to initialize.
//
// If successful, Listen will run a go routine that listen for DNS requests, and call the appropriate registered callbacks.
// The go routine can be closed via Close, and any packet decoding errors are sent to errChan
func (d *DNSListener) Listen(dataChan chan Packet) error {
	d.active = true

	dcap, err := dnscap.New()
	if err != nil {
		return err
	}

	go (func() {
		for {
			dnsReqs := dcap.Read()

			fmt.Println(len(dnsReqs))

			for dnsReq := range dnsReqs {
				fmt.Println(dnsReq.Request)

				for regex, f := range d.handlers {
					if dnsReq.Error != nil {
						dataChan <- Packet{
							Error: dnsReq.Error,
							Host:  "",
						}

						continue
					}

					dataChan <- Packet{
						Error: nil,
						Host:  dnsReq.Request,
					}

					if regex.MatchString(dnsReq.Request) {
						f(d, dnsReq.Request)
					}

					if !d.active {
						dcap.Close()
						return
					}
				}
			}

			fmt.Println("onto the next")
		}
	})()

	return nil
}
