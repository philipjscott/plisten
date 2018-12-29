package dnscap

import "testing"

func TestGenerateDNSFilter(t *testing.T) {
	ips := []string{"172.16.254.1", "172.36.224.1", "192.168.1.15"}
	expected := "udp and port 53 and (src host 172.16.254.1 or src host 172.36.224.1 or src host 192.168.1.15)"

	if generateDNSFilter(ips) != expected {
		t.Fail()
	}
}
