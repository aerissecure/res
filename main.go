package main

import (
	"net"
	"fmt"
	"flag"
	"os"
	"strings"
)

func main() {
	ipv4 := flag.Bool("4", false, "ipv4 addresses only")
	ipv6 := flag.Bool("6", false, "ipv6 addresses only")
	l := flag.Bool("l", false, "list ips")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("no hostnames specified")
		os.Exit(1)
	}
	if *ipv4 && *ipv6 {
		fmt.Println("use -4 or -6, not both")
		os.Exit(1)
	}

	hosts := flag.Args()
	lookups := make(chan map[string][]net.IP)

	for _, h := range hosts {
		go func(h string) {
			m := make(map[string][]net.IP)
			m[h] = []net.IP{}
			ips, _ := net.LookupHost(h)
			for _, ip := range ips {
				m[h] = append(m[h], net.ParseIP(ip))
			}
			lookups <- m
		}(h)
	}

	hostMap := make(map[string][]net.IP)
	for i := 0; i < len(hosts); i++ {
		m := <- lookups
		for h, ips := range m {
			hostMap[h] = ips
		}
	}

	if *l {
		out := []net.IP{}
		for _, ips := range hostMap {
			for _, ip := range ips {
				if !*ipv4 && !*ipv6 {
					out = append(out, ip)
				} else if *ipv4 && ip.To4() != nil {
					out = append(out, ip)
				} else if *ipv6 && ip.To4() == nil {
					out = append(out, ip)
				}
			}
		}
		o := []string{}
		for _, ip := range out {
			o = append(o, ip.String())
		}
		fmt.Println(strings.Join(o, " "))
	} else {
		for h, ips := range hostMap {
			fmt.Println(h)
			for _, ip := range ips {
				if !*ipv4 && !*ipv6 {
					fmt.Println("  ", ip.String())
				} else if *ipv4 && ip.To4() != nil {
					fmt.Println("  ", ip.String())
				} else if *ipv6 && ip.To4() == nil {
					fmt.Println("  ", ip.String())
				}
			}
			fmt.Println()
		}
	}
}