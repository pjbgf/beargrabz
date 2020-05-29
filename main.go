package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	maxPayloadLines uint16        = 10
	device          string        = "any"
	snapLen         int32         = 1024
	promisc         bool          = true
	timeout         time.Duration = 15 * time.Second
)

func main() {
	handle, err := pcap.OpenLive(device, snapLen, promisc, timeout)
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for pkt := range source.Packets() {
		printBearerTokens(pkt)
	}
}

func printBearerTokens(packet gopacket.Packet) {
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		payload := string(applicationLayer.Payload())
		if strings.Contains(payload, "HTTP") {
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if ipLayer == nil || tcpLayer == nil {
				return
			}

			ip, _ := ipLayer.(*layers.IPv4)
			tcp, _ := tcpLayer.(*layers.TCP)
			if ip == nil || tcp == nil {
				return
			}

			r := regexp.MustCompile(`^(GET|POST|Host:|Authorization:\sBearer)\s+(\S+)`)
			scanner := bufio.NewScanner(strings.NewReader(payload))
			scanner.Split(bufio.ScanLines)

			bearerFound := false
			headerLines := ""
			for scanner.Scan() && maxPayloadLines > 0 {
				headers := r.FindStringSubmatch(scanner.Text())
				if len(headers) > 0 {
					if !bearerFound {
						bearerFound = strings.HasPrefix(headers[0], "Authorization")
						if bearerFound {
							headerLines += fmt.Sprintf("%v*********\n", headers[0][:30])
							continue
						}
					}
					headerLines += fmt.Sprintf("%v\n", headers[0])
				}
			}

			// Only print if some headers were found and one of them is the bearer token
			if len(headerLines) > 0 && bearerFound {
				fmt.Printf("%s:%d -> %s:%d\n", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort)
				fmt.Printf("%v\n", headerLines)
			}
		}
	}
}
