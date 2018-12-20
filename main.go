package main

import (
    "fmt"
    "log"
    "github.com/google/gopacket/pcap"
)

func isValidDevice(device pcap.Interface) bool {
    const loopbackMask = 0x01
    const runningMask = 0x04

    if (device.Flags & loopbackMask != 0 || device.Flags & runningMask == 0) {
        return false
    }

    return len(device.Addresses) > 0
}  

func validDevices(devices []pcap.Interface) {
    // Print device information
    fmt.Println("Devices found:")

    for _, device := range devices {
        if !isValidDevice(device) {
            continue
        }

        fmt.Println("\nName: ", device.Name)
        fmt.Println("Description: ", device.Description)
        fmt.Println("Devices addresses: ", device.Description)
        for _, address := range device.Addresses {
            fmt.Println("- IP address: ", address.IP)
            fmt.Println("- Subnet mask: ", address.Netmask)
        }
    }
}

func main() {
    // Find all devices
    devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatal(err)
    }

    validDevices(devices)
}