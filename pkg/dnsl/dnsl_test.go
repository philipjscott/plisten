package dnsl

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestDNSListen(t *testing.T) {
	// TODO: Implement TUN/TAP interface for testing; it won't with travis-ci otherwise!
	return

	dataChan := make(chan Packet)
	dl := New()
	url := "http://www.facebook.com"
	calledCallback := false

	err := dl.Register(".*", func(d *DNSListener, match string) {
		t.Log("Callback called!")
		calledCallback = true

		if !strings.Contains(url, match) {
			t.Fail()
		}
		d.Close()
	})

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = dl.Listen(dataChan)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = http.Get(url)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Sleep two seconds, since there is a 1 second timeout buffer for dnsl
	time.Sleep(time.Second * 2)

	// Do a non-blocking check to see if an error occured
	select {
	case data := <-dataChan:
		if data.Error != nil {
			t.Log(data.Error)
			t.Fail()
		}
	default:
		t.Log("No packet errors occurred")
	}

	if !calledCallback {
		t.Fail()
		t.Log("Failed to call callback")
	}
}
