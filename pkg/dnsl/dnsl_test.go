package dnsl

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestDNSListen(t *testing.T) {
	var errChan chan error
	dl := New()
	url := "http://www.facebook.com"
	calledCallback := false

	err := dl.Register(".*", func(d *DNSListener, match string) {
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

	err = dl.Listen(errChan)
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
	time.Sleep(time.Second * 3)

	// Check if packet error occurred
	err = <-errChan
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if !calledCallback {
		t.Fail()
		t.Log("Failed to call callback")
	}
}
