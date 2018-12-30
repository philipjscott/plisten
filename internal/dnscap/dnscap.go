package dnscap

import (
	"errors"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/google/gopacket/pcap"
)

type DNSCapture struct {
	Request string
	Error   error
}

type DNSCapturer struct {
	handles []*pcap.Handle
	parser  *gopacket.DecodingLayerParser
	decoded []gopacket.LayerType
	dns     *layers.DNS
}

func New() (*DNSCapturer, error) {
	var (
		eth     layers.Ethernet
		ip4     layers.IPv4
		ip6     layers.IPv6
		tcp     layers.TCP
		udp     layers.UDP
		dns     layers.DNS
		payload gopacket.Payload
	)

	devices, err := getActiveDevices()
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, errors.New("No active devices found!")
	}

	handles, err := getHandles(devices)
	if err != nil {
		return nil, err
	}

	err = setDNSFilters(handles, devices)
	if err != nil {
		return nil, err
	}

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &payload)
	decoded := make([]gopacket.LayerType, 0, 10)

	return &DNSCapturer{
		handles,
		parser,
		decoded,
		&dns,
	}, nil
}

func (d *DNSCapturer) Read() chan DNSCapture {
	var chans []chan gopacket.Packet
	agg := make(chan DNSCapture)

	for _, handle := range d.handles {
		chans = append(chans, gopacket.NewPacketSource(handle, handle.LinkType()).Packets())
	}

	for _, ch := range chans {
		go func(c chan gopacket.Packet) {
			for packet := range c {
				err := d.parser.DecodeLayers(packet.Data(), &d.decoded)

				if err != nil {
					agg <- DNSCapture{"", err}
				}

				for _, typ := range d.decoded {
					switch typ {
					case layers.LayerTypeDNS:
						for _, dnsQuestion := range d.dns.Questions {
							agg <- DNSCapture{string(dnsQuestion.Name), nil}
						}
					}
				}

			}
		}(ch)
	}

	return agg
}

func (d *DNSCapturer) Close() {
	for _, handle := range d.handles {
		handle.Close()
	}
}

func getHandle(deviceName string) (*pcap.Handle, error) {
	const (
		snapshotLen int32         = 1024
		promiscuous bool          = false
		timeout     time.Duration = time.Second * 1
	)

	return pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
}

func getHandles(devices []pcap.Interface) ([]*pcap.Handle, error) {
	var handles []*pcap.Handle

	for _, device := range devices {
		handle, err := getHandle(device.Name)

		if err != nil {
			return nil, err
		}

		handles = append(handles, handle)
	}

	return handles, nil
}

func getActiveDevices() ([]pcap.Interface, error) {
	var activeDevices []pcap.Interface
	devices, err := pcap.FindAllDevs()

	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if isActiveDevice(device) {
			activeDevices = append(activeDevices, device)
		}
	}

	return activeDevices, nil
}

func isActiveDevice(device pcap.Interface) bool {
	const loopbackMask = 0x01
	const runningMask = 0x04

	if device.Flags&loopbackMask != 0 || device.Flags&runningMask == 0 {
		return false
	}

	return len(device.Addresses) > 0
}

func setDNSFilters(handles []*pcap.Handle, devices []pcap.Interface) error {
	var ips []string

	for _, device := range devices {
		for _, address := range device.Addresses {
			ips = append(ips, address.IP.String())
		}
	}

	filter := generateDNSFilter(ips)

	// TODO: Make a different filter for each device
	for _, handle := range handles {
		err := handle.SetBPFFilter(filter)

		if err != nil {
			return err
		}
	}

	return nil
}

func generateDNSFilter(ips []string) string {
	filter := "udp and port 53 and ("

	for i, ip := range ips {
		if i != 0 {
			filter += " or "
		}

		filter += "src host " + ip
	}

	filter += ")"

	return filter
}
