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

func setFilter(handle *pcap.Handle, ipAddr string) error {
    return handle.SetBPFFilter("tcp")
}

func main() {
    var tcp layers.TCP
    const (
        snapshotLen int32  = 1024
        promiscuous  bool   = false
        timeout      time.Duration = 30 * time.Second
    )

    device, err := getDevice()
    
    if err != nil {
        log.Fatal(err)
    }

    ipAddr := device.Addresses[0].IP.String()
    handle, err := pcap.OpenLive(device.Name, snapshotLen, promiscuous, timeout)

    defer handle.Close()

    if err != nil {
        log.Fatal(err)
    }

    err = setFilter(handle, ipAddr)

    if err != nil {
        log.Fatal(err)
    }

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    
    for packet := range packetSource.Packets() {
        parser := gopacket.NewDecodingLayerParser(layers.LayerTypeTCP, &tcp)
        foundLayerTypes := []gopacket.LayerType{}

        err := parser.DecodeLayers(packet.Data(), &foundLayerTypes)
        if err != nil {
            fmt.Println("Trouble decoding layers: ", err)
        }

        fmt.Println(string(tcp.LayerContents()))
    }
}