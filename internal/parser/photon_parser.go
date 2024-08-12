package sniffer

import (
	"log"

	"github.com/google/gopacket/pcap"
)

func Start() error {
	handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	return nil
}
