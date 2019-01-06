// The dnsl package listens for packets and decodes DNS information;
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
	"regexp"

	"github.com/ScottyFillups/plisten/internal/dnscap"
)

type DNSListener struct {
	handlers map[*regexp.Regexp]func(*DNSListener, string)
	active   bool
}

// Returns a DNSListener struct.
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

// Stop listening for packets
func (d *DNSListener) Close() {
	d.active = false
}

// Listen begins listening for packets on all active internet connections (promiscuous mode is disabled).
// Returns an error if the packet capturer fails to initialize.
//
// If successful, Listen will run a go routine that listen for DNS requests, and call the appropriate registered callbacks.
// The go routine can be closed via Close, and any packet decoding errors are sent to errChan
func (d *DNSListener) Listen(errChan chan error) error {
	d.active = true

	dcap, err := dnscap.New()
	if err != nil {
		return err
	}

	go (func() {
		for {
			dnsReqs := dcap.Read()

			for dnsReq := range dnsReqs {
				for regex, f := range d.handlers {
					if dnsReq.Error != nil {
						errChan <- dnsReq.Error
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
		}
	})()

	return nil
}
