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

func New() DNSListener {
	return DNSListener{
		handlers: make(map[*regexp.Regexp]func(*DNSListener, string)),
		active:   false,
	}
}

func (d *DNSListener) Register(regexStr string, handler func(*DNSListener, string)) error {
	regex := regexp.MustCompile(regexStr)

	if regex == nil {
		return errors.New("Failed to compile regular expression")
	}

	d.handlers[regex] = handler

	return nil
}

func (d *DNSListener) Close() {
	d.active = false
}

func (d *DNSListener) Listen() error {
	d.active = true

	dcap, err := dnscap.New()
	if err != nil {
		return err
	}
	defer dcap.Close()

	for {
		dnsReqs := dcap.Read()
		
		for dnsReq := range dnsReqs {
			for regex, f := range d.handlers {
				if dnsReq.Error != nil {
					return dnsReq.Error
				}

				if regex.MatchString(dnsReq.Request) {
					f(d, dnsReq.Request)
				}

				if !d.active {
					return nil
				}
			}
		}
	}
}
