package main

import (
    "fmt"
    "log"
    "time"
    "errors"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket"
)

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

    if (device.Flags & loopbackMask != 0 || device.Flags & runningMask == 0) {
        return false
    }

    return len(device.Addresses) > 0
}

func setDNSFilter(handle *pcap.Handle, ipAddr string) error {
    return handle.SetBPFFilter("udp and port 53 and src host " + ipAddr)
}

func main() {
    var (
        eth layers.Ethernet
        ip4 layers.IPv4
        ip6 layers.IPv6
        tcp layers.TCP
        udp layers.UDP
        dns layers.DNS
        payload gopacket.Payload
    )
    const (
        snapshotLen int32  = 1024
        promiscuous  bool   = false
        timeout      time.Duration = 30 * time.Second
        filter  string = "udp and port 53 and src host "
    )

    device, err := getDevice()
    
    if err != nil {
        log.Fatal(err)
    }

    ipAddr := device.Addresses[0].IP.String()
    
    handle, err := pcap.OpenLive(device.Name, snapshotLen, promiscuous, pcap.BlockForever)

    defer handle.Close()

    if err != nil {
        log.Fatal(err)
    }

    err = setDNSFilter(handle, ipAddr)

    if err != nil {
        log.Fatal(err)
    }

    parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &payload)
    decodedLayers := make([]gopacket.LayerType, 0, 10)

	for {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			fmt.Println("Error reading packet data: ", err)
			continue
		}

		err = parser.DecodeLayers(data, &decodedLayers)
		for _, typ := range decodedLayers {
			switch typ {
			case layers.LayerTypeDNS:
                for _, dnsQuestion := range dns.Questions {
                    timeFormatted := time.Now().Format("2006-01-02 15:04:05")
                    fmt.Println(timeFormatted, "->", string(dnsQuestion.Name))
                }
			}
		}

		if err != nil {
			fmt.Println("Error encountered: ", err)
		}
	}
}