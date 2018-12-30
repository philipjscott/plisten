package dnsl

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestDNSListen(t *testing.T) {
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

	err = dl.Listen(nil)
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
	if !calledCallback {
		t.Fail()
		t.Log("Failed to call callback")
	}
}
