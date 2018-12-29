package dnsl

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDNSListen(t *testing.T) {
	dl := New()
	url := "www.facebook.com"

	err := dl.Register(".*", func(d *DNSListener, match string) {
		fmt.Println("foo")
		if url != match {
			t.Fail()
		}
		d.Close()
	})

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = dl.Listen()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = http.Get(url)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
