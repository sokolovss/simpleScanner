package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	defaultPort = "80"
)

var (
	waitGroup sync.WaitGroup
)

// retrieves the IP address of the local machine.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// scans the local network subnet for open TCP ports on each IP address
func main() {
	subnet := getSubnet()
	for i := 1; i <= 255; i++ {
		waitGroup.Add(1)
		go func(j int) {
			defer waitGroup.Done()
			address := fmt.Sprintf("%s.%d:%s", subnet, j, defaultPort)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				if strings.Contains(err.Error(), "too many open files") {
					fmt.Println("Increase ulimit.")
				}
				return
			}
			conn.Close()
			fmt.Printf("%s has open port %s\n", address, defaultPort)
		}(i)
	}
	waitGroup.Wait()
}

// returns the subnet of the local IP address by removing the last octet.
func getSubnet() string {
	ipAddress := getLocalIP()
	ipParts := strings.Split(ipAddress, ".")
	subnet := strings.Join(ipParts[:len(ipParts)-1], ".")
	return subnet
}
