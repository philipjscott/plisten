package dnscap

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"regexp"
)

type DNSCapturer struct {
	handlers map[*regexp.Regexp]func(*DNSCapturer, string)
	active   bool
}

func New() DNSCapturer {
	return DNSCapturer{
		handlers: make(map[*regexp.Regexp]func(*DNSCapturer, string)),
		active:   false,
	}
}

func (d *DNSCapturer) Register(regexStr string, handler func(*DNSCapturer, string)) error {
	regex := regexp.MustCompile(regexStr)

	if regex == nil {
		return errors.New("Failed to compile regular expression")
	}

	d.handlers[regex] = handler

	return nil
}

func (d *DNSCapturer) Close() {
	d.active = false
}

func (d *DNSCapturer) Listen() {
	d.active = true

	const (
		snapshotLen int32  = 1024
		promiscuous bool   = false
		filter      string = "udp and port 53 and src host "
	)
	var (
		eth     layers.Ethernet
		ip4     layers.IPv4
		ip6     layers.IPv6
		tcp     layers.TCP
		udp     layers.UDP
		dns     layers.DNS
		payload gopacket.Payload
	)

	device, err := getDevice()
	if err != nil {
		log.Fatal(err)
	}

	handle, err := pcap.OpenLive(device.Name, snapshotLen, promiscuous, pcap.BlockForever)
	defer handle.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = setDNSFilter(handle, device)
	if err != nil {
		log.Fatal(err)
	}

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &payload)
	decodedLayers := make([]gopacket.LayerType, 0, 10)

	for {
		if !d.active {
			break
		}

		data, _, err := handle.ReadPacketData()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading packet data: ", err)
			continue
		}

		err = parser.DecodeLayers(data, &decodedLayers)
		for _, typ := range decodedLayers {
			switch typ {
			case layers.LayerTypeDNS:
				for _, dnsQuestion := range dns.Questions {
					for regex, f := range d.handlers {
						dnsStr := string(dnsQuestion.Name)
						if regex.MatchString(dnsStr) {
							f(d, dnsStr)
						}
					}
				}
			}
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error encountered: ", err)
		}
	}
}

func getHandle(deviceName string) (*pcap.Handle, error) {
	const (
		snapshotLen int32 = 1024
		promiscuous bool  = false
	)

	return pcap.OpenLive(deviceName, snapshotLen, promiscuous, pcap.BlockForever)
}

func getDevice() (*pcap.Interface, error) {
	devices, err := pcap.FindAllDevs()

	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if isActiveDevice(device) {
			return &device, nil
		}
	}

	return nil, errors.New("Failed to find active device")
}

func isActiveDevice(device pcap.Interface) bool {
	const loopbackMask = 0x01
	const runningMask = 0x04

	if device.Flags&loopbackMask != 0 || device.Flags&runningMask == 0 {
		return false
	}

	return len(device.Addresses) > 0
}

func setDNSFilter(handle *pcap.Handle, device *pcap.Interface) error {
	filter := "udp and port 53 and ("

	for i, address := range device.Addresses {
		if i != 0 {
			filter += "or "
		}

		filter += "src host " + address.IP.String()
	}

	filter += ")"

	return handle.SetBPFFilter(filter)
}
