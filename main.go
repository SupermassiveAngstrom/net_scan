package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	userIP, err := getUserIP()
	if err != nil {
		fmt.Println("Error getting user IP:", err)
		return
	}

	defaultStartIP := getNetworkStartIP(userIP)
	var startIP, endIP string
	fmt.Printf("User IP: %s\n", userIP)
	fmt.Printf("Default Starting IP: %s\n", defaultStartIP)
	fmt.Print("Enter the starting IP address (or press Enter to use default): ")
	fmt.Scanln(&startIP)

	if startIP == "" {
		startIP = defaultStartIP
	}

	fmt.Print("Enter the ending IP address (e.g., 192.168.1.254): ")
	fmt.Scanln(&endIP)

	startIP = strings.TrimSpace(startIP)
	endIP = strings.TrimSpace(endIP)

	startIPAddr := net.ParseIP(startIP)
	endIPAddr := net.ParseIP(endIP)

	if startIPAddr == nil || endIPAddr == nil {
		fmt.Println("Invalid IP address format.")
		return
	}

	fmt.Println("Scanning range from", startIP, "to", endIP)

	start := ipToInt(startIPAddr)
	end := ipToInt(endIPAddr)

	for ip := start; ip <= end; ip++ {
		ipAddr := intToIP(ip)
		go pingIP(ipAddr)
	}

	time.Sleep(5 * time.Second)
}

func getUserIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no valid IP address found")
}

func getNetworkStartIP(userIP string) string {
	ipParts := strings.Split(userIP, ".")
	ipParts[3] = "1"
	return strings.Join(ipParts, ".")
}

func ipToInt(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func intToIP(ipInt uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ipInt>>24), byte(ipInt>>16&0xFF), byte(ipInt>>8&0xFF), byte(ipInt&0xFF))
}

func pingIP(ip string) {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)

	err := cmd.Run()
	if err == nil {
		fmt.Printf("IP %s is up\n", ip)
	}
}
