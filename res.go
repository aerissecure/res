/*
functionality:
- perform forward lookup returning all ips for host
- perform reverse lookup returning all hosts for ip
- recursive mode where the host/ip looking is performed on the result of the initial request
- return output in json, -oJ
eventually:
- perform arin queries
- perform whois queries
- let recursive take a number for the number of times to recurse?

*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

type Addr struct {
	Addr    string  `json:"addr"`
	Lookups []*Addr `json:"lookups"`
	// Type     string  `json:"type"` // host, ip
	resolver *net.Resolver
}

func (i *Addr) typeIP() bool {
	if ip := net.ParseIP(i.Addr); ip != nil {
		return true
	}
	return false
}

func (i *Addr) typeHost() bool {
	return !i.typeIP()
}

// may need to remove duplicates
// func flatten(addrs []*Addr) []*Addr {
// 	var flat []*Addr
// 	for _, addr := range addrs {
// 		if len(addr.Lookups) > 0 {
// 			// flat = append(flat, addr) need to blank out the lookups in this case
// 			flat = append(flat, flatten(addr.Lookups)...)
// 		} else {
// 			flat = append(flat, addr)
// 		}
// 	}
// 	return flat
// }

func (i *Addr) resolve(ipv4, ipv6 bool) {
	if i.typeIP() {
		names, err := i.resolver.LookupAddr(context.Background(), i.Addr)
		// typecheck error ???
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
		for _, n := range names {
			i.Lookups = append(i.Lookups, &Addr{Addr: n, resolver: i.resolver})
		}
	}
	if i.typeHost() {
		ips, err := net.LookupHost(i.Addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
		for _, ip := range ips {
			typedIP := net.ParseIP(ip)
			if typedIP == nil {
				continue
			}
			isIPv4 := !(typedIP.To4() == nil)
			isIPv6 := typedIP != nil && !isIPv4
			if (ipv4 && isIPv4) || (ipv6 && isIPv6) {
				i.Lookups = append(i.Lookups, &Addr{Addr: ip, resolver: i.resolver})
			}
		}
	}
}

// automatically detect Addr as hostname or ip
func main() {
	var ipv4, ipv6, recursive, flat, pretty bool
	flag.BoolVar(&ipv4, "4", true, "include ipv4 responses")
	flag.BoolVar(&ipv6, "6", false, "include ipv6 responses")
	flag.BoolVar(&recursive, "r", false, "recursive (resolve lookups)")
	flag.BoolVar(&flat, "f", false, "flatten output") // json only?
	flag.BoolVar(&pretty, "p", false, "pretty print")
	_ = flat // not sure what flat really means right now...
	_ = pretty

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "error: provide at least one ip/hostname as an argument")
		os.Exit(1)
	}

	resolver := &net.Resolver{PreferGo: true}

	var addrs []*Addr
	for _, arg := range flag.Args() {
		addrs = append(addrs, &Addr{Addr: arg, resolver: resolver})
	}

	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr *Addr) {
			defer wg.Done()
			addr.resolve(ipv4, ipv6)
			if recursive {
				var wgr sync.WaitGroup
				for _, addr := range addr.Lookups {
					wgr.Add(1)
					go func(addr *Addr) {
						defer wgr.Done()
						addr.resolve(ipv4, ipv6)
					}(addr)
				}
				wgr.Wait()
			}
		}(addr)
	}
	wg.Wait()

	// for _, addr := range addrs {
	// 	fmt.Println(addr.Addr)
	// 	for _, addr := range addr.Lookups {
	// 		fmt.Printf("\t%s\n", addr.Addr)
	// 		if recursive {
	// 			for _, addr := range addr.Lookups {
	// 				fmt.Printf("\t\t%s\n", addr.Addr)
	// 			}
	// 		}
	// 	}
	// }

	// if flat {
	// 	addrs = flatten(addrs)
	// }

	var b []byte
	if pretty {
		b, _ = json.MarshalIndent(addrs, "", "  ")
	} else {
		b, _ = json.Marshal(addrs)
	}
	fmt.Printf("%s\n", string(b))
}
