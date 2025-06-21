/* 
by Andrew Mamdouh
github: https://github.com/Pwn3rx0
linkedin: https://www.linkedin.com/in/andrew-mamdouh122/
*/ 

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func url_to_domain(domain string) string {
	if strings.HasPrefix(domain, "http://") {
		domain = strings.TrimPrefix(domain, "http://")
	} else if strings.HasPrefix(domain, "https://") {
		domain = strings.TrimPrefix(domain, "https://")
	}
	return domain
}

func saveIPsToFile(ips []string, file *os.File) {
	if file == nil {
		return
	}
	for _, ip := range ips {
		file.WriteString(ip + "\n")
	}
}

func resolveDomain(domain string, ipv4Out, ipv6Out *os.File) {
	domain = url_to_domain(domain)
	ips, err := net.LookupHost(domain)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var ipv4List, ipv6List []string

	for _, ip := range ips {
		if net.ParseIP(ip).To4() != nil {
			ipv4List = append(ipv4List, ip)
		} else {
			ipv6List = append(ipv6List, ip)
		}
	}

	if ipv4Out != nil || (ipv4Out == nil && ipv6Out == nil) {
		if ipv4Out != nil {
			saveIPsToFile(ipv4List, ipv4Out)
			if len(ipv4List) > 0 {
				fmt.Println("[+] IPv4 IPs saved in file")
			}
		} else {
			fmt.Printf("IPv4 of %v:\n", domain)
			for _, ip := range ipv4List {
				fmt.Println(ip)
			}
		}
	}

	if ipv6Out != nil || (ipv4Out == nil && ipv6Out == nil) {
		if ipv6Out != nil {
			saveIPsToFile(ipv6List, ipv6Out)
			if len(ipv6List) > 0 {
				fmt.Println("[+] IPv6 IPs saved in file")
			}
		} else {
			fmt.Printf("IPv6 of %v:\n", domain)
			for _, ip := range ipv6List {
				fmt.Println(ip)
			}
		}
	}
}



func main() {
	domain := flag.String("d", "", "target domain to resolve")
	file := flag.String("f", "", "file containing domains")
	ipv4path := flag.String("ipv4", "", "Save IPv4 results to file")
	ipv6path := flag.String("ipv6", "", "Save IPv6 results to file")
	flag.Parse()

	var ipv4Out, ipv6Out *os.File
	var err error

	if *ipv4path != "" {
		ipv4Out, err = os.Create(*ipv4path)
		if err != nil {
			fmt.Println("Failed to create IPv4 output file:", err)
			return
		}
		defer ipv4Out.Close()
	}

	if *ipv6path != "" {
		ipv6Out, err = os.Create(*ipv6path)
		if err != nil {
			fmt.Println("Failed to create IPv6 output file:", err)
			return
		}
		defer ipv6Out.Close()
	}

	if *domain != "" {
		resolveDomain(*domain, ipv4Out, ipv6Out)
	} else if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Println("Failed to open file:", err)
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				resolveDomain(line, ipv4Out, ipv6Out)
			}
		}
	} else {
		fmt.Println("Usage: ipresolver -d <domain> OR -f <file> [-ipv4 <file>] [-ipv6 <file>]")
	}
}
